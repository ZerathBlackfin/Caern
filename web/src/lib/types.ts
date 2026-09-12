export interface Tile {
  name: string
  url: string
  icon?: string
  description?: string
  width: number
  height: number
  newTab: boolean
}

export interface Block {
  section?: string
  tiles: Tile[]
}

export interface Page {
  title: string
  theme: string
  backgroundColor?: string
  backgroundImage?: string
  columns: number
  blocks: Block[]
}

export interface Problem {
  file: string
  line?: number
  message: string
}

export interface State {
  page: Page
  problems?: Problem[]
}
