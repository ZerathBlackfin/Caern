import type { State } from './types'

export const live = $state<{ state: State | null; connected: boolean }>({
  state: null,
  connected: true,
})

export function connect() {
  const events = new EventSource('/api/events')
  events.addEventListener('state', (e) => {
    live.state = JSON.parse(e.data)
    live.connected = true
  })
  events.addEventListener('error', () => {
    live.connected = false
  })
  return () => events.close()
}
