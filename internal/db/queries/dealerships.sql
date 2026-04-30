-- name: ListDealerships :many
SELECT id, name, open_time, close_time, created_at, updated_at
FROM dealerships
ORDER BY id;

-- name: GetDealershipByID :one
SELECT id, name, open_time, close_time, created_at, updated_at
FROM dealerships
WHERE id = $1
LIMIT 1;

-- name: SearchDealershipsByName :many
SELECT id, name, open_time, close_time, created_at, updated_at
FROM dealerships
WHERE unaccent(name) ILIKE unaccent('%' || $1 || '%')
ORDER BY name;

-- name: UpdateDealershipByID :one
UPDATE dealerships
SET
  name = COALESCE(sqlc.narg('name')::text, name),
  open_time = COALESCE(sqlc.narg('open_time')::time, open_time),
  close_time = COALESCE(sqlc.narg('close_time')::time, close_time),
  updated_at = now()
WHERE id = $1
RETURNING id, name, open_time, close_time, created_at, updated_at;

-- name: DeleteDealershipByID :exec
DELETE FROM dealerships
WHERE id = $1;

-- name: CreateDealership :one
INSERT INTO dealerships (name, open_time, close_time)
VALUES ($1, $2, $3)
RETURNING id, name, open_time, close_time, created_at, updated_at;
