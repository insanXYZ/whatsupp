CREATE TABLE IF NOT EXISTS message_attachments (
  id SERIAL NOT NULL,
  message_id INT NOT NULL,
  file_url TEXT NOT NULL,
  file_ext VARCHAR(15) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id),
  FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE
);
