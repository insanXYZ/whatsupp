import "dotenv/config";

type WsEventFunc = ((this: WebSocket, ev: Event) => any) | null;
type WsMessageEventFunc = ((this: WebSocket, ev: MessageEvent) => any) | null;
type WsCloseFunc = ((this: WebSocket, ev: CloseEvent) => any) | null;

type ConnectWSProps = {
  onOpen: WsEventFunc;
  onError: WsEventFunc;
  onMessage: WsMessageEventFunc;
  onClose: WsCloseFunc;
};

export const ConnectWS = ({
  onOpen,
  onError,
  onMessage,
  onClose,
}: ConnectWSProps): WebSocket => {
  const ws = new WebSocket(process.env.NEXT_PUBLIC_BASE_URL_WS!);
  ws.onopen = onOpen;

  ws.onmessage = onMessage;

  ws.onclose = onClose;

  ws.onerror = onError;

  return ws;
};
