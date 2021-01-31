interface InitialData {
  payload: string;
}

export let ws: WebSocket
export let savedInitialData: InitialData

export enum WSReadyState {
  CONNECTING,
  OPEN,
  CLOSING,
  CLOSED
}

export function createWebSocket (token: string, onReady: (status: string) => void): typeof ws {
  const websocket = new WebSocket('wss://stop.parkfig.com/ws')

  websocket.addEventListener('close', () => {
    onReady('close')
  })
  websocket.addEventListener('error', () => {
    onReady('error')
  })
  websocket.onopen = () => {
    ws = websocket

    ws.send(token)

    onReady('open')

    ws.addEventListener('message', (e) => {
      const res = JSON.parse(e.data)
      
      // TODO: sync res to store
      console.log(res);
    })
  }

  return ws
}