CREATE TABLE IF NOT EXISTS message_attachments (
  id VARCHAR(100) NOT NULL,
  message_id VARCHAR(100) NOT NULL,
  file_url TEXT NOT NULL,
  file_type VARCHAR(15) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(message_id) REFERENCES messages(id)
);
