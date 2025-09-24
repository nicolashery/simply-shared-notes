CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE spaces (
  id INTEGER PRIMARY KEY,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  admin_token TEXT NOT NULL,
  edit_token TEXT NOT NULL,
  view_token TEXT NOT NULL
, created_by INTEGER REFERENCES members(id) ON DELETE SET NULL, updated_by INTEGER REFERENCES members(id) ON DELETE SET NULL);
CREATE INDEX idx_spaces_admin_token ON spaces(admin_token);
CREATE INDEX idx_spaces_edit_token ON spaces(edit_token);
CREATE INDEX idx_spaces_view_token ON spaces(view_token);
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
CREATE TABLE activity (
  id INTEGER PRIMARY KEY,
  created_at DATETIME NOT NULL,
  space_id INTEGER NOT NULL,
  public_id TEXT NOT NULL,
  member_id INTEGER,
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id INTEGER,

  FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
  FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE SET NULL
);
CREATE INDEX idx_activities_space_id_created_at ON activity(space_id, created_at DESC);
CREATE INDEX idx_activities_entity_type_entity_id ON activity(entity_type, entity_id);
CREATE UNIQUE INDEX idx_activities_space_id_public_id ON activity(space_id, public_id);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20250409141508'),
  ('20250420184520'),
  ('20250822150859'),
  ('20250924151138');
