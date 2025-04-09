CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE spaces (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  admin_token TEXT NOT NULL
);
CREATE INDEX idx_spaces_admin_token ON spaces(admin_token);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20250409141508');
