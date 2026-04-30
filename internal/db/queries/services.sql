-- SERVICES

-- name: CreateService :one
INSERT INTO services (required_bay_type_id, name, anticipated_minutes)
VALUES ($1, $2, $3)
RETURNING id, required_bay_type_id, name, anticipated_minutes, created_at, updated_at;

-- name: UpdateServiceByID :one
UPDATE services
SET
  required_bay_type_id = COALESCE(sqlc.narg('required_bay_type_id')::int, required_bay_type_id),
  name = COALESCE(sqlc.narg('name')::text, name),
  anticipated_minutes = COALESCE(sqlc.narg('anticipated_minutes')::int, anticipated_minutes),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, required_bay_type_id, name, anticipated_minutes, created_at, updated_at;

-- name: DeleteServiceByID :exec
DELETE FROM services
WHERE id = $1;

-- name: ListServices :many
SELECT id, required_bay_type_id, name, anticipated_minutes, created_at, updated_at
FROM services
ORDER BY name;

-- name: GetServiceByID :one
SELECT id, required_bay_type_id, name, anticipated_minutes, created_at, updated_at
FROM services
WHERE id = $1
LIMIT 1;

-- name: SearchServicesByName :many
SELECT id, required_bay_type_id, name, anticipated_minutes, created_at, updated_at
FROM services
WHERE unaccent(name) ILIKE unaccent('%' || $1 || '%')
ORDER BY name;

-- SERVICE REQUIREMENTS (skills required by services)

-- name: AddSkillRequirementToService :exec
INSERT INTO service_requirements (service_id, skill_id)
VALUES ($1, $2)
ON CONFLICT (service_id, skill_id) DO NOTHING;

-- name: AddSkillRequirementsToService :exec
INSERT INTO service_requirements (service_id, skill_id)
SELECT sqlc.arg('service_id'), unnest(sqlc.arg('skill_ids')::int[])
ON CONFLICT (service_id, skill_id) DO NOTHING;

-- name: RemoveSkillRequirementFromService :exec
DELETE FROM service_requirements
WHERE service_id = $1
  AND skill_id = $2;

-- name: RemoveSkillRequirementsFromService :exec
DELETE FROM service_requirements
WHERE service_id = sqlc.arg('service_id')
  AND skill_id = ANY(sqlc.arg('skill_ids')::int[]);

-- name: ListSkillIDsByServiceID :many
SELECT skill_id
FROM service_requirements
WHERE service_id = $1
ORDER BY skill_id;

-- name: GetServiceDetailByID :one
-- Includes required skill_ids and required bay type info
SELECT
  s.id,
  s.required_bay_type_id,
  sbt.name AS required_bay_type_name,
  s.name,
  s.anticipated_minutes,
  s.created_at,
  s.updated_at,
  COALESCE(
    jsonb_agg(sr.skill_id) FILTER (WHERE sr.skill_id IS NOT NULL),
    '[]'::jsonb
  ) AS required_skill_ids
FROM services s
JOIN service_bay_types sbt ON sbt.id = s.required_bay_type_id
LEFT JOIN service_requirements sr ON sr.service_id = s.id
WHERE s.id = $1
GROUP BY s.id, sbt.id
LIMIT 1;


