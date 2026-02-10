import "dotenv/config";

export const ConnectWS = (): WebSocket => {
  const ws = new WebSocket(process.env.NEXT_PUBLIC_BASE_URL_WS!);
  return ws;
};
