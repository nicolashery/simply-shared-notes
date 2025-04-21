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
