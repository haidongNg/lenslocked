CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  sys_user_id INT UNIQUE,
  token_hash TEXT UNIQUE NOT NULL
);