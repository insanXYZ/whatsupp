CREATE TABLE IF NOT EXISTS message_attachments (
  id SERIAL NOT NULL,
  message_id INT NOT NULL,
  file_url TEXT NOT NULL,
  file_ext VARCHAR(15) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(message_id) REFERENCES messages(id)
);

CREATE TRIGGER trigger_set_created_at_message_attachment
BEFORE INSERT ON message_attachments
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

