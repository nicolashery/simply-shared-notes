-- name: CreateNote :one
INSERT INTO notes (
  created_at,
  updated_at,
  created_by,
  updated_by,
  space_id,
  public_id,
  title,
  content
) VALUES (
  @created_at,
  @updated_at,
  @created_by,
  @updated_by,
  @space_id,
  @public_id,
  @title,
  @content
) RETURNING *;

-- name: GetNoteByPublicID :one
SELECT * FROM notes
WHERE space_id = @space_id
  AND public_id = @public_id
LIMIT 1;

-- name: ListNotes :many
SELECT
    n.public_id,
    n.created_at,
    cb.public_id as created_by_public_id,
    cb.name as created_by_name,
    n.updated_at,
    ub.public_id as updated_by_public_id,
    ub.name as updated_by_name,
    n.title
FROM
    notes n
    LEFT JOIN members cb ON cb.id = n.created_by
    LEFT JOIN members ub ON ub.id = n.updated_by
WHERE n.space_id = @space_id
ORDER BY
    n.updated_at DESC,
    n.id DESC;

-- name: UpdateNote :one
UPDATE notes
SET updated_at = @updated_at,
  updated_by = @updated_by,
  title = @title,
  content = @content
WHERE id = @note_id
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = ?;
