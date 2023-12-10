CREATE TABLE sys_users (
    id INT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    passwword_hash TEXT NOT NULL
);