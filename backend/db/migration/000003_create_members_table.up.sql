CREATE OR REPLACE FUNCTION set_joined_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.joined_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'type_group') THEN
      CREATE TYPE role_group_member AS ENUM ('ADMIN','MEMBER');
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS members (
  id SERIAL NOT NULL,
  group_id INT NOT NULL,
  user_id INT NOT NULL,
  role role_group_member NOT NULL,
  joined_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(group_id) REFERENCES groups(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);


CREATE TRIGGER trigger_set_joined_at_member
BEFORE INSERT ON members
FOR EACH ROW
EXECUTE FUNCTION set_joined_at();

