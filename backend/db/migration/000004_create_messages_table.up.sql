CREATE TABLE IF NOT EXISTS messages (
  id SERIAL NOT NULL,
  conversation_id INT NOT NULL,
  user_id INT NOT NULL,
  message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,
  PRIMARY KEY(id),
  FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

DO
$$BEGIN
    CREATE TRIGGER trigger_set_updated_at_message
    BEFORE UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
EXCEPTION
   WHEN duplicate_object THEN
      NULL;
END;$$;


