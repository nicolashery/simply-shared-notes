-- name: CreateSpace :one
INSERT INTO spaces (
  created_at,
  updated_at,
  name,
  email,
  admin_token,
  edit_token,
  view_token
) VALUES (
  @created_at,
  @updated_at,
  @name,
  @email,
  @admin_token,
  @edit_token,
  @view_token
) RETURNING *;

-- name: GetSpaceByAccessToken :one
SELECT * FROM spaces
WHERE admin_token = @token
  OR edit_token = @token
  OR view_token = @token
LIMIT 1;

-- name: GetSpaceStats :one
SELECT
  COUNT(DISTINCT n.id) AS notes_count,
  COUNT(DISTINCT m.id) AS members_count
FROM spaces s
LEFT JOIN notes n ON n.space_id = s.id
LEFT JOIN members m ON m.space_id = s.id
WHERE s.id = @space_id;

-- name: UpdateSpaceCreatedBy :exec
UPDATE spaces
SET created_by = @created_by,
  updated_by = @created_by
WHERE id = @space_id;

-- name: UpdateSpace :one
UPDATE spaces
SET updated_at = @updated_at,
  updated_by = @updated_by,
  name = @name
WHERE id = @space_id
RETURNING *;
