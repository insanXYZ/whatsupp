CREATE TABLE IF NOT EXISTS messages (
  id INT NOT NULL,
  member_id INT NOT NULL,
  message TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ,
  PRIMARY KEY(id),
  FOREIGN KEY(member_id) REFERENCES members(id)
);

CREATE TRIGGER trigger_set_created_at_message
BEFORE INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

CREATE TRIGGER trigger_set_updated_at_message
BEFORE UPDATE ON messages
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

