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
