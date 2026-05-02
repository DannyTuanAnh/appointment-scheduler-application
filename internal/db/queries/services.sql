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

-- name: DeleteServiceByID :execrows
DELETE FROM services
WHERE id = $1;

-- name: ListServices :many
SELECT s.id, s.required_bay_type_id, sbt.name as type_name, s.name as service_name, s.anticipated_minutes, s.created_at, s.updated_at
FROM services s
LEFT JOIN service_bay_types sbt ON sbt.id = s.required_bay_type_id
ORDER BY s.name;

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
  (
    SELECT COALESCE(jsonb_agg(sk.name), '[]'::jsonb)
    FROM service_requirements sr
    JOIN skills sk ON sk.id = sr.skill_id
    WHERE sr.service_id = s.id
  ) AS required_skills_name
FROM services s
JOIN service_bay_types sbt ON sbt.id = s.required_bay_type_id
WHERE s.id = $1;

-- name: GetSkillRequirementIDs :many
SELECT skill_id
FROM service_requirements
WHERE service_id = $1;

-- name: SearchServicesByName :many
SELECT s.id, s.required_bay_type_id, sbt.name as type_name, s.name as service_name, s.anticipated_minutes, s.created_at, s.updated_at
FROM services s
LEFT JOIN service_bay_types sbt ON sbt.id = s.required_bay_type_id
WHERE unaccent(s.name) ILIKE unaccent('%' || sqlc.arg('name')::text || '%')
ORDER BY s.name;

-- SERVICE REQUIREMENTS (skills required by services)

-- name: AddSkillRequirementsToService :exec
INSERT INTO service_requirements (service_id, skill_id)
SELECT sqlc.arg('service_id'), unnest(sqlc.arg('skill_ids')::int[])
ON CONFLICT (service_id, skill_id) DO NOTHING;

-- name: RemoveSkillRequirementsFromService :execrows
DELETE FROM service_requirements
WHERE service_id = sqlc.arg('service_id')
  AND skill_id = ANY(sqlc.arg('skill_ids')::int[]);





