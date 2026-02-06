import 'dotenv/config'

export const ConnectWS = (): WebSocket => {
    const ws = new WebSocket(process.env.BASE_URL_WS!)
    return ws
}