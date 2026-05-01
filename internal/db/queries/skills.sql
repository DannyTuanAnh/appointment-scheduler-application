-- name: CreateSkill :one
INSERT INTO skills (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: ListSkills :many
SELECT id, name, created_at, updated_at
FROM skills
ORDER BY name;

-- name: GetSkillByID :one
SELECT id, name, created_at, updated_at
FROM skills
WHERE id = $1
LIMIT 1;

-- name: SearchSkillsByName :many
SELECT id, name, created_at, updated_at
FROM skills
WHERE unaccent(name) ILIKE unaccent('%' || $1 || '%')
ORDER BY name;

-- name: UpdateSkillNameByID :one
UPDATE skills
SET name = COALESCE(sqlc.narg('name')::text, name),
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, name, created_at, updated_at;

-- name: DeleteSkillByID :execrows
DELETE FROM skills
WHERE id = $1;
