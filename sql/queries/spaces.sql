-- name: GetSpaceByAccessToken :one
SELECT * FROM spaces
WHERE admin_token = ? LIMIT 1;
