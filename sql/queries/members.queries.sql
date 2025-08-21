-- name: CreateMember :one
INSERT INTO members (
  created_at,
  updated_at,
  created_by,
  updated_by,
  space_id,
  public_id,
  name
) VALUES (
  @created_at,
  @updated_at,
  @created_by,
  @updated_by,
  @space_id,
  @public_id,
  @name
) RETURNING *;

-- name: UpdateMemberCreatedBy :exec
UPDATE members
SET created_by = @created_by,
  updated_by = @created_by
WHERE id = @member_id;

-- name: GetMemberByID :one
SELECT * FROM members
WHERE id = @id
LIMIT 1;

-- name: GetMemberByPublicID :one
SELECT * FROM members
WHERE space_id = @space_id
  AND public_id = @public_id
LIMIT 1;

-- name: ListMembers :many
SELECT * FROM members
WHERE space_id = @space_id
ORDER BY name;

-- name: UpdateMember :one
UPDATE members
SET updated_at = @updated_at,
  updated_by = @updated_by,
  name = @name
WHERE id = @member_id
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM members WHERE id = ?;
