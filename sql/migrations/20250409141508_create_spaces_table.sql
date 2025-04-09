-- migrate:up
CREATE TABLE spaces (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  admin_token TEXT NOT NULL
);

CREATE INDEX idx_spaces_admin_token ON spaces(admin_token);

-- migrate:down

