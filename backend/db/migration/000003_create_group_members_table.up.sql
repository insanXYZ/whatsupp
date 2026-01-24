CREATE TYPE role_group_member AS ENUM ('ADMIN','MEMBER')

CREATE TABLE IF NOT EXISTS group_members (
  id VARCHAR(100) NOT NULL,
  group_id VARCHAR(100) NOT NULL,
  user_id VARCHAR(100) NOT NULL,
  role role_group_member NOT NULL,
  joined_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY(group_id) REFERENCES groups(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
