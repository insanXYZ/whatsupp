CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN 
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_created_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.created_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(100) NOT NULL,
  image TEXT NOT NULL,
  bio VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ,
  PRIMARY KEY(id)
);

CREATE TRIGGER trigger_set_created_at_user
BEFORE INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

CREATE TRIGGER trigger_set_updated_at_user
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();



