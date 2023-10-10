CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS player_group (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id uuid,
  name VARCHAR(50) NOT NULL,
  public boolean,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);