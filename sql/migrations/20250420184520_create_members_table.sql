-- migrate:up

CREATE TABLE members (
  id INTEGER PRIMARY KEY,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  created_by INTEGER,
  updated_by INTEGER,
  space_id INTEGER NOT NULL,
  public_id TEXT NOT NULL,
  name TEXT NOT NULL,

  FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES members(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES members(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX idx_members_space_id_public_id ON members(space_id, public_id);

ALTER TABLE spaces ADD COLUMN created_by INTEGER REFERENCES members(id) ON DELETE SET NULL;
ALTER TABLE spaces ADD COLUMN updated_by INTEGER REFERENCES members(id) ON DELETE SET NULL;

-- migrate:down
