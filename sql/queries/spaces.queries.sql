-- name: CreateSpace :one
INSERT INTO spaces (
  created_at,
  updated_at,
  name,
  email,
  admin_token,
  edit_token,
  view_token
) VALUES ( ?, ?, ?, ?, ?, ?, ? ) RETURNING *;

-- name: GetSpaceByAccessToken :one
SELECT * FROM spaces
WHERE admin_token = ? LIMIT 1;
