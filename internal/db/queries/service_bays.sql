-- name: ListServiceBays :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
ORDER BY sb.dealership_id, sb.bay_type_id, sb.name;

-- name: ListServiceBaysByDealershipID :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.dealership_id = sqlc.arg('dealership_id')
ORDER BY sb.bay_type_id, sb.name;

-- name: ListServiceBaysByTypeID :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.bay_type_id = sqlc.arg('bay_type_id')
ORDER BY sb.dealership_id, sb.name;

-- name: ListServiceBaysByDealershipIDAndTypeID :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.dealership_id = sqlc.arg('dealership_id')
  AND sb.bay_type_id = sqlc.arg('bay_type_id')
ORDER BY sb.id;

-- name: GetServiceBayByID :one
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.id = sqlc.arg('id')
LIMIT 1;

-- name: SearchServiceBaysByName :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE unaccent(sb.name) ILIKE unaccent('%' || sqlc.arg('service_bay_name')::text || '%')
ORDER BY sb.dealership_id, sb.name;

-- name: SearchServiceBaysByNameAndDealershipID :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.dealership_id = sqlc.arg('dealership_id')
  AND unaccent(sb.name) ILIKE unaccent('%' || sqlc.arg('service_bay_name')::text || '%')
ORDER BY sb.bay_type_id, sb.name;

-- name: SearchServiceBaysByNameAndTypeID :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.bay_type_id = sqlc.arg('bay_type_id')
  AND unaccent(sb.name) ILIKE unaccent('%' || sqlc.arg('service_bay_name')::text || '%')
ORDER BY sb.dealership_id, sb.name;

-- name: SearchServiceBaysByNameDealershipIDAndTypeID :many
SELECT sb.id, sb.dealership_id, sb.bay_type_id, d.name as dealership_name, sbt.name as type_name, sb.name as service_bay_name, sb.is_active, sb.created_at, sb.updated_at
FROM service_bays sb
JOIN dealerships d ON sb.dealership_id = d.id
LEFT JOIN service_bay_types sbt ON sb.bay_type_id = sbt.id
WHERE sb.dealership_id = sqlc.arg('dealership_id')
  AND sb.bay_type_id = sqlc.arg('bay_type_id')
  AND unaccent(sb.name) ILIKE unaccent('%' || sqlc.arg('service_bay_name')::text || '%')
ORDER BY sb.name;

-- name: CreateServiceBay :one
INSERT INTO service_bays (dealership_id, bay_type_id, name)
VALUES ($1, $2, $3)
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

-- name: DeleteServiceBayByID :execrows
DELETE FROM service_bays
WHERE id = $1;


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

-- name: DeleteServiceBayTypeByID :execrows
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
