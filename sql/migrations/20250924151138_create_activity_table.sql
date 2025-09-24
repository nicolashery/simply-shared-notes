-- migrate:up
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
CREATE INDEX idx_activities_member_id ON activity(member_id);
CREATE INDEX idx_activities_entity_type_entity_id ON activity(entity_type, entity_id);
CREATE UNIQUE INDEX idx_activities_space_id_public_id ON activity(space_id, public_id);

-- migrate:down
