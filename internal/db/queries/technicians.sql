-- name: CreateTechnician :one
INSERT INTO technicians (dealership_id, name, level)
VALUES ($1, $2, $3)
RETURNING id, dealership_id, name, level, is_active, inactive_since, created_at, updated_at;

-- name: SetTechnicianOnLeave :one
UPDATE technicians
SET is_active = FALSE,
    inactive_since = now(),
    updated_at = now()
WHERE id = $1
RETURNING id, dealership_id, name, level, is_active, inactive_since, created_at, updated_at;

-- name: SetTechnicianBackToWork :one
UPDATE technicians
SET is_active = TRUE,
    inactive_since = NULL,
    updated_at = now()
WHERE id = $1
RETURNING id, dealership_id, name, level, is_active, inactive_since, created_at, updated_at;

-- name: TransferTechnicianDealership :one
UPDATE technicians
SET dealership_id = $2,
    updated_at = now()
WHERE id = $1
RETURNING id, dealership_id, name, level, is_active, inactive_since, created_at, updated_at;

-- name: UpdateTechnicianInfoByID :one
UPDATE technicians
SET
  name = COALESCE(sqlc.narg('name')::text, name),
  level = COALESCE(sqlc.narg('level')::technician_level, level),
  updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, dealership_id, name, level, is_active, inactive_since, created_at, updated_at;

-- name: DeleteTechnicianIfInactiveOverOneMonth :exec
DELETE FROM technicians
WHERE id = $1
  AND is_active = FALSE
  AND inactive_since IS NOT NULL
  AND inactive_since < (now() - interval '1 month');

-- name: DeleteTechnicianByID :execrows
DELETE FROM technicians
WHERE id = $1;

-- name: AddSkillsToTechnician :exec
INSERT INTO technician_skills (technician_id, skill_id)
SELECT sqlc.arg('technician_id'), unnest(sqlc.arg('skill_ids')::int[])
ON CONFLICT (technician_id, skill_id) DO NOTHING;

-- name: RemoveSkillsFromTechnician :execrows
DELETE FROM technician_skills
WHERE technician_id = sqlc.arg('technician_id')
  AND skill_id = ANY(sqlc.arg('skill_ids')::int[]);

-- name: ListTechniciansByDealershipID :many
SELECT t.id, t.dealership_id, d.name as dealership_name, t.name as technician_name, t.level, t.is_active, t.inactive_since, t.created_at, t.updated_at
FROM technicians t
JOIN dealerships d ON t.dealership_id = d.id
WHERE t.dealership_id = $1
ORDER BY t.dealership_id, t.name;

-- name: SearchTechniciansByName :many
SELECT t.id, t.dealership_id, d.name as dealership_name, t.name as technician_name, t.level, t.is_active, t.inactive_since, t.created_at, t.updated_at
FROM technicians t
JOIN dealerships d ON t.dealership_id = d.id
WHERE unaccent(t.name) ILIKE unaccent('%' || sqlc.arg('technician_name')::text || '%')
ORDER BY t.dealership_id, t.name;

-- name: SearchTechniciansByNameAndDealershipID :many
SELECT t.id, t.dealership_id, d.name as dealership_name, t.name as technician_name, t.level, t.is_active, t.inactive_since, t.created_at, t.updated_at
FROM technicians t
JOIN dealerships d ON t.dealership_id = d.id
WHERE dealership_id = $1
  AND unaccent(t.name) ILIKE unaccent('%' || sqlc.arg('technician_name')::text || '%')
ORDER BY t.dealership_id, t.name;

-- name: FindActiveTechniciansByDealershipWithRequiredSkills :many
-- Returns active technician IDs who have ALL of the required skills (skill_ids).
-- If skill_ids is NULL/empty, returns all active technician IDs for the dealership.
SELECT t.id AS technician_id
FROM technicians t
WHERE t.dealership_id = sqlc.arg('dealership_id')
  AND t.is_active = TRUE
  AND (
    sqlc.arg('skill_ids')::int[] IS NULL
    OR cardinality(sqlc.arg('skill_ids')::int[]) = 0
    OR (
      SELECT array_agg(DISTINCT ts.skill_id)
      FROM technician_skills ts
      WHERE ts.technician_id = t.id
    ) @> sqlc.arg('skill_ids')::int[]
  )
ORDER BY t.id;

-- name: GetDetailTechnicianByID :one
SELECT
  t.id, t.dealership_id, d.name as dealership_name, t.name as technician_name,
  t.level, t.is_active, t.inactive_since, t.created_at, t.updated_at,
  (
    SELECT COALESCE(jsonb_agg(s.name), '[]'::jsonb)
    FROM technician_skills ts
    JOIN skills s ON s.id = ts.skill_id
    WHERE ts.technician_id = t.id
  ) AS skills
FROM technicians t
JOIN dealerships d ON t.dealership_id = d.id
WHERE t.id = $1;

-- name: GetTechnicianByID :one
SELECT t.id, t.dealership_id, d.name as dealership_name, t.name as technician_name, t.level, t.is_active, t.inactive_since, t.created_at, t.updated_at
FROM technicians t
JOIN dealerships d ON t.dealership_id = d.id
WHERE t.id = $1;
