import { live } from './live.svelte'
import type { Block } from './types'

export type Resize = {
  id: string
  pxWidth: number
  pxHeight: number
  width: number
  height: number
}

export type Drag = {
  id: string
  left: number
  top: number
  width: number
  height: number
}

export const edit = $state<{
  on: boolean
  blocks: Block[]
  drag: Drag | null
  resize: Resize | null
  error: string | null
}>({
  on: false,
  blocks: [],
  drag: null,
  resize: null,
  error: null,
})

export function start() {
  edit.blocks = $state.snapshot(live.state?.page.blocks ?? []) as Block[]
  edit.on = true
  edit.error = null
}

export function stop() {
  document.documentElement.classList.remove('grabbing', 'resizing')
  edit.on = false
  edit.blocks = []
  edit.drag = null
  edit.resize = null
  edit.error = null
}

export async function save() {
  const blocks = edit.blocks.map((block) => ({
    id: block.id ?? '',
    tiles: block.tiles.map((tile) => ({ id: tile.id, size: `${tile.width}x${tile.height}` })),
  }))
  const response = await fetch('/api/layout', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ blocks }),
  })
  if (response.ok) {
    stop()
    return
  }
  edit.error = (await response.text()).trim() || `the server answered ${response.status}`
}

export function startDrag(id: string, event: PointerEvent) {
  event.preventDefault()
  event.stopPropagation()

  const tile = (event.target as HTMLElement).closest<HTMLElement>('[data-tile]')
  if (!tile) return
  const box = tile.getBoundingClientRect()
  const grabX = event.clientX - box.left
  const grabY = event.clientY - box.top
  edit.drag = { id, left: box.left, top: box.top, width: box.width, height: box.height }
  document.documentElement.classList.add('grabbing')

  let settled = 0

  const move = (e: PointerEvent) => {
    if (!edit.drag) return
    edit.drag.left = e.clientX - grabX
    edit.drag.top = e.clientY - grabY
    if (performance.now() < settled) return
    if (dragOver(id, e.clientX, e.clientY)) settled = performance.now() + SETTLE_MS
  }
  const up = () => {
    edit.drag = null
    document.documentElement.classList.remove('grabbing')
    window.removeEventListener('pointermove', move)
    window.removeEventListener('pointerup', up)
  }
  window.addEventListener('pointermove', move)
  window.addEventListener('pointerup', up)
}

export function startResize(id: string, event: PointerEvent) {
  event.preventDefault()
  event.stopPropagation()

  const tile = (event.target as HTMLElement).closest<HTMLElement>('[data-tile]')
  const grid = tile?.closest<HTMLElement>('[data-block]')
  const found = locate(id)
  if (!tile || !grid || !found) return

  const style = getComputedStyle(grid)
  const gap = parseFloat(style.columnGap) || 12
  const columns = style.gridTemplateColumns.split(' ').length
  const cell = (grid.getBoundingClientRect().width - gap * (columns - 1)) / columns
  const row = parseFloat(style.gridAutoRows) || 58
  const tiles = edit.blocks[found.block].tiles
  const box = tile.getBoundingClientRect()

  const span = (px: number, step: number, max: number) => clamp(Math.round((px + gap) / step), 1, max)
  const cellStep = cell + gap
  const rowStep = row + gap

  document.documentElement.classList.add('resizing')
  edit.resize = {
    id,
    pxWidth: box.width,
    pxHeight: box.height,
    width: tiles[found.index].width,
    height: tiles[found.index].height,
  }

  const move = (e: PointerEvent) => {
    const at = locate(id)
    if (!edit.resize || !at) return

    const from = tile.getBoundingClientRect()
    const pxWidth = clamp(e.clientX - from.left, cell, columns * cellStep - gap)
    const pxHeight = clamp(e.clientY - from.top, row, MAX_SPAN * rowStep - gap)
    edit.resize.pxWidth = pxWidth
    edit.resize.pxHeight = pxHeight
    edit.resize.width = span(pxWidth, cellStep, columns)
    edit.resize.height = span(pxHeight, rowStep, MAX_SPAN)

    const target = edit.blocks[at.block].tiles[at.index]
    target.width = edit.resize.width
    target.height = edit.resize.height
  }
  const up = () => {
    window.removeEventListener('pointermove', move)
    window.removeEventListener('pointerup', up)
    document.documentElement.classList.remove('resizing')
    edit.resize = null
  }
  window.addEventListener('pointermove', move)
  window.addEventListener('pointerup', up)
}

const MAX_SPAN = 12
const SETTLE_MS = 160

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max)
}

function locate(id: string) {
  for (const [block, { tiles }] of edit.blocks.entries()) {
    const index = tiles.findIndex((tile) => tile.id === id)
    if (index >= 0) return { block, index }
  }
  return null
}

function dragOver(id: string, x: number, y: number) {
  const from = locate(id)
  const under = document.elementFromPoint(x, y)
  const grid = under?.closest<HTMLElement>('[data-block]')
  if (!from || !grid) return false

  const block = Number(grid.dataset.block)
  const over = under?.closest<HTMLElement>('[data-tile]')
  if (over?.dataset.tile === id) return false

  const tiles = edit.blocks[block].tiles
  if (!over && tiles.length > 0) return false

  let index = 0
  if (over) {
    const box = over.getBoundingClientRect()
    index = tiles.findIndex((tile) => tile.id === over.dataset.tile)
    if (x > box.left + box.width / 2) index++
  }
  if (block === from.block && from.index < index) index--
  if (block === from.block && index === from.index) return false

  const [tile] = edit.blocks[from.block].tiles.splice(from.index, 1)
  edit.blocks[block].tiles.splice(index, 0, tile)
  return true
}
