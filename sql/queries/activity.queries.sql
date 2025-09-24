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
