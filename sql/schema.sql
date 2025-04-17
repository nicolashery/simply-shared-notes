CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE spaces (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  admin_token TEXT NOT NULL,
  edit_token TEXT NOT NULL,
  view_token TEXT NOT NULL
);
CREATE INDEX idx_spaces_admin_token ON spaces(admin_token);
CREATE INDEX idx_spaces_edit_token ON spaces(edit_token);
CREATE INDEX idx_spaces_view_token ON spaces(view_token);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20250409141508');
