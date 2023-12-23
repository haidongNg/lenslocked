CREATE TABLE sys_sessions (
  id SERIAL PRIMARY KEY,
  sys_user_id INT UNIQUE PEFERENCES sys_users(id) ON DELETE CASCADE,
  token_hash TEXT UNIQUE NOT NULL
);