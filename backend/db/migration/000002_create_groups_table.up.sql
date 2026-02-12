DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_group') THEN
      CREATE TYPE type_group AS ENUM ( 'GROUP', 'PERSONAL');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS groups (
  id SERIAL NOT NULL,
  name VARCHAR(100) NOT NULL,
  bio VARCHAR(100) NOT NULL,
  type type_group NOT NULL,
  image TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,
  PRIMARY KEY(id)
);

CREATE TRIGGER trigger_set_updated_at_group
BEFORE UPDATE ON groups
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

