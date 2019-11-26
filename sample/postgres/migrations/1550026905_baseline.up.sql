CREATE TABLE IF NOT EXISTS users (
  id         INTEGER                                   NOT NULL PRIMARY KEY,
  login      VARCHAR(100)                              NOT NULL,
  password   VARCHAR(100)                              NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
  is_super   BOOLEAN DEFAULT FALSE                     NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS users_id_uindex
  ON users (id);

CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex
  ON users (login);
