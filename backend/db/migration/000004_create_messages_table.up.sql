CREATE TABLE IF NOT EXISTS messages (
  id SERIAL NOT NULL,
  member_id INT NOT NULL,
  message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,
  PRIMARY KEY(id),
  FOREIGN KEY(member_id) REFERENCES members(id) ON DELETE CASCADE
);

CREATE TRIGGER trigger_set_updated_at_message
BEFORE UPDATE ON messages
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

