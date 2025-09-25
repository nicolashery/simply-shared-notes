-- name: CreateActivity :one
INSERT INTO activity (
  created_at,
  space_id,
  public_id,
  member_id,
  action,
  entity_type,
  entity_id
) VALUES (
  @created_at,
  @space_id,
  @public_id,
  @member_id,
  @action,
  @entity_type,
  @entity_id
) RETURNING *;

-- name: SetActivityEntityIDToNull :exec
UPDATE activity
SET entity_id = NULL
WHERE entity_type = @entity_type
  AND entity_id = @entity_id;

-- name: ListActivity :many
SELECT * FROM activity
WHERE space_id = @space_id
ORDER BY
    created_at DESC,
    id DESC
LIMIT @limit;

-- name: GetActivityByPublicID :one
SELECT * FROM activity
WHERE space_id = @space_id
  AND public_id = @public_id
LIMIT 1;
