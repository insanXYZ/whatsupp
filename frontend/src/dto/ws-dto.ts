interface SendMessageRequest {
  message: string;
  group_id: number | undefined;
  receiver_id: number | undefined;
}
