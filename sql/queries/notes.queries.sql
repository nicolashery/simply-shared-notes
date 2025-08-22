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
SELECT * FROM notes
WHERE space_id = @space_id
ORDER BY title;

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
