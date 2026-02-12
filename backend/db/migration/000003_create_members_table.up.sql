DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_group_member') THEN
      CREATE TYPE role_group_member AS ENUM ('ADMIN','MEMBER');
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS members (
  id SERIAL NOT NULL,
  group_id INT NOT NULL,
  user_id INT NOT NULL,
  role role_group_member NOT NULL,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id),
  UNIQUE(group_id, user_id),
  FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TRIGGER trigger_set_joined_at_member
BEFORE INSERT ON members
FOR EACH ROW
EXECUTE FUNCTION set_joined_at();

