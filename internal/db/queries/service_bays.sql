-- name: ListServiceBays :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
ORDER BY dealership_id, id;

-- name: GetServiceBayByID :one
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE id = $1
LIMIT 1;

-- name: SearchServiceBaysByName :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE unaccent(name) ILIKE unaccent('%' || $1 || '%')
ORDER BY name;

-- name: ListServiceBaysByDealershipID :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE dealership_id = $1
ORDER BY id;

-- name: SearchServiceBaysByNameAndDealershipID :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE dealership_id = $1
  AND unaccent(name) ILIKE unaccent('%' || $2 || '%')
ORDER BY name;

-- name: CreateServiceBay :one
INSERT INTO service_bays (dealership_id, bay_type_id, name, is_active)
VALUES ($1, $2, $3, COALESCE($4, TRUE))
RETURNING id, dealership_id, bay_type_id, name, is_active, created_at, updated_at;

-- name: UpdateServiceBayByID :one
UPDATE service_bays
SET
  dealership_id = COALESCE(sqlc.narg('dealership_id')::int, dealership_id),
  bay_type_id = COALESCE(sqlc.narg('bay_type_id')::int, bay_type_id),
  name = COALESCE(sqlc.narg('name')::text, name),
  is_active = COALESCE(sqlc.narg('is_active')::boolean, is_active),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, dealership_id, bay_type_id, name, is_active, created_at, updated_at;

-- name: DeleteServiceBayByID :exec
DELETE FROM service_bays
WHERE id = $1;

-- name: ListServiceBaysByTypeID :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE bay_type_id = $1
ORDER BY dealership_id, id;

-- name: SearchServiceBaysByNameAndTypeID :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE bay_type_id = $1
  AND unaccent(name) ILIKE unaccent('%' || $2 || '%')
ORDER BY name;

-- name: SearchServiceBaysByNameDealershipIDAndTypeID :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE dealership_id = $1
  AND bay_type_id = $2
  AND unaccent(name) ILIKE unaccent('%' || $3 || '%')
ORDER BY name;

-- name: ListServiceBaysByDealershipIDAndTypeID :many
SELECT id, dealership_id, bay_type_id, name, is_active, created_at, updated_at
FROM service_bays
WHERE dealership_id = $1
  AND bay_type_id = $2
ORDER BY id;

-- SERVICE BAY TYPES

-- name: CreateServiceBayType :one
INSERT INTO service_bay_types (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: UpdateServiceBayTypeByID :one
UPDATE service_bay_types
SET name = COALESCE(sqlc.narg('name')::text, name),
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, name, created_at, updated_at;

-- name: DeleteServiceBayTypeByID :exec
DELETE FROM service_bay_types
WHERE id = $1;

-- name: ListServiceBayTypes :many
SELECT id, name, created_at, updated_at
FROM service_bay_types
ORDER BY name;

-- name: GetServiceBayTypeByID :one
SELECT id, name, created_at, updated_at
FROM service_bay_types
WHERE id = $1
LIMIT 1;

-- name: SearchServiceBayTypesByName :many
SELECT id, name, created_at, updated_at
FROM service_bay_types
WHERE unaccent(name) ILIKE unaccent('%' || $1 || '%')
ORDER BY name;
