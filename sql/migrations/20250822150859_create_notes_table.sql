-- migrate:up

CREATE TABLE notes (
  id INTEGER PRIMARY KEY,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  created_by INTEGER,
  updated_by INTEGER,
  space_id INTEGER NOT NULL,
  public_id TEXT NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,

  FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES members(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES members(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX idx_notes_space_id_public_id ON notes(space_id, public_id);

-- migrate:down
