DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_conversations') THEN
      CREATE TYPE type_conversations AS ENUM ( 'GROUP', 'PERSONAL');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS conversations (
  id SERIAL NOT NULL,
  name VARCHAR(100) NOT NULL,
  bio VARCHAR(100) NOT NULL,
  type type_conversations NOT NULL,
  image TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,
  PRIMARY KEY(id)
);

DO
$$BEGIN
    CREATE TRIGGER trigger_set_updated_at_conversation
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
EXCEPTION
   WHEN duplicate_object THEN
      NULL;
END;$$;


