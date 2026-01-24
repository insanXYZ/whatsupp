CREATE TABLE IF NOT EXISTS messages (
  id VARCHAR(100) NOT NULL,
  group_id VARCHAR(100) NOT NULL,
  sender_id VARCHAR(100) NOT NULL,
  message TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(group_id) REFERENCES groups(id),
  FOREIGN KEY(sender_id) REFERENCES users(id)
);
