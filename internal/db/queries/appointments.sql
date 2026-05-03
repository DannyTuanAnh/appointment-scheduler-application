-- name: GetAppointmentsOfBayOrTechnicianInTimeRange :many
WITH target_bays AS (
    SELECT id
    FROM service_bays sb
    WHERE sb.dealership_id = sqlc.arg('dealership_id')
        AND sb.bay_type_id = sqlc.arg('bay_type_id')
        AND sb.is_active = TRUE
),
target_technicians AS (
    SELECT t.id
    FROM technicians t

    JOIN technician_skills ts ON t.id = ts.technician_id

    WHERE t.dealership_id = sqlc.arg('dealership_id')
    AND t.is_active = TRUE

    AND ts.skill_id = ANY(sqlc.arg('skill_ids')::int[])
    GROUP BY t.id

    HAVING COUNT(DISTINCT ts.skill_id) = cardinality(sqlc.arg('skill_ids')::int[])
)

SELECT 
    ap.bay_id, 
    ap.technician_id, 
    ap.duration 
FROM appointments ap
WHERE status NOT IN ('cancelled', 'no_show') 
  AND (ap.bay_id IN (SELECT id FROM target_bays)  OR ap.technician_id IN (SELECT id FROM target_technicians))
  AND ap.duration && tstzrange(sqlc.arg('from_time')::timestamptz, sqlc.arg('to_time')::timestamptz, '[)')
ORDER BY lower(ap.duration);

-- Common projection with joined info
-- Note: duration is tstzrange; use lower(duration) as start time, upper(duration) as end time.

-- name: ListAppointments :many
SELECT
  a.id,
  a.dealership_id,
  d.name AS dealership_name,
  d.open_time AS dealership_open_time,
  d.close_time AS dealership_close_time,
  a.service_id,
  sv.name AS service_name,
  sv.anticipated_minutes,
  a.bay_id,
  sb.name AS bay_name,
  a.technician_id,
  t.name AS technician_name,
  a.customer_name,
  a.status,
  a.duration,
  lower(a.duration) AS start_time,
  upper(a.duration) AS end_time,
  a.created_at,
  a.updated_at
FROM appointments a
JOIN dealerships d ON d.id = a.dealership_id
JOIN services sv ON sv.id = a.service_id
LEFT JOIN technicians t ON t.id = a.technician_id
LEFT JOIN service_bays sb ON sb.id = a.bay_id
ORDER BY lower(a.duration) DESC, a.id DESC;

-- name: ListAppointmentsByDealershipInTimeRange :many
SELECT
  a.id,
  a.dealership_id,
  d.name AS dealership_name,
  d.open_time AS dealership_open_time,
  d.close_time AS dealership_close_time,
  a.service_id,
  sv.name AS service_name,
  sv.anticipated_minutes,
  a.bay_id,
  sb.name AS bay_name,
  a.technician_id,
  t.name AS technician_name,
  a.customer_name,
  a.status,
  a.duration,
  lower(a.duration) AS start_time,
  upper(a.duration) AS end_time,
  a.created_at,
  a.updated_at
FROM appointments a
JOIN dealerships d ON d.id = a.dealership_id
JOIN services sv ON sv.id = a.service_id
LEFT JOIN technicians t ON t.id = a.technician_id
LEFT JOIN service_bays sb ON sb.id = a.bay_id
WHERE a.dealership_id = sqlc.arg('dealership_id')
  AND a.duration && tstzrange(sqlc.arg('from_time'), sqlc.arg('to_time'), '[)')
ORDER BY lower(a.duration), a.id;

-- name: ListAppointmentsByTechnicianInTimeRange :many
SELECT
  a.id,
  a.dealership_id,
  d.name AS dealership_name,
  d.open_time AS dealership_open_time,
  d.close_time AS dealership_close_time,
  a.service_id,
  sv.name AS service_name,
  sv.anticipated_minutes,
  a.bay_id,
  sb.name AS bay_name,
  a.technician_id,
  t.name AS technician_name,
  a.customer_name,
  a.status,
  a.duration,
  lower(a.duration) AS start_time,
  upper(a.duration) AS end_time,
  a.created_at,
  a.updated_at
FROM appointments a
JOIN dealerships d ON d.id = a.dealership_id
JOIN services sv ON sv.id = a.service_id
LEFT JOIN technicians t ON t.id = a.technician_id
LEFT JOIN service_bays sb ON sb.id = a.bay_id
WHERE a.technician_id = sqlc.arg('technician_id')
  AND a.duration && tstzrange(sqlc.arg('from_time'), sqlc.arg('to_time'), '[)')
ORDER BY lower(a.duration), a.id;

-- name: ListAppointmentsByServiceBayInTimeRange :many
SELECT
  a.id,
  a.dealership_id,
  d.name AS dealership_name,
  d.open_time AS dealership_open_time,
  d.close_time AS dealership_close_time,
  a.service_id,
  sv.name AS service_name,
  sv.anticipated_minutes,
  a.bay_id,
  sb.name AS bay_name,
  a.technician_id,
  t.name AS technician_name,
  a.customer_name,
  a.status,
  a.duration,
  lower(a.duration) AS start_time,
  upper(a.duration) AS end_time,
  a.created_at,
  a.updated_at
FROM appointments a
JOIN dealerships d ON d.id = a.dealership_id
JOIN services sv ON sv.id = a.service_id
LEFT JOIN technicians t ON t.id = a.technician_id
LEFT JOIN service_bays sb ON sb.id = a.bay_id
WHERE a.bay_id = sqlc.arg('bay_id')
  AND a.duration && tstzrange(sqlc.arg('from_time'), sqlc.arg('to_time'), '[)')
ORDER BY lower(a.duration), a.id;

-- name: ListAppointmentsByServiceInTimeRange :many
SELECT
  a.id,
  a.dealership_id,
  d.name AS dealership_name,
  d.open_time AS dealership_open_time,
  d.close_time AS dealership_close_time,
  a.service_id,
  sv.name AS service_name,
  sv.anticipated_minutes,
  a.bay_id,
  sb.name AS bay_name,
  a.technician_id,
  t.name AS technician_name,
  a.customer_name,
  a.status,
  a.duration,
  lower(a.duration) AS start_time,
  upper(a.duration) AS end_time,
  a.created_at,
  a.updated_at
FROM appointments a
JOIN dealerships d ON d.id = a.dealership_id
JOIN services sv ON sv.id = a.service_id
LEFT JOIN technicians t ON t.id = a.technician_id
LEFT JOIN service_bays sb ON sb.id = a.bay_id
WHERE a.service_id = sqlc.arg('service_id')
  AND a.duration && tstzrange(sqlc.arg('from_time'), sqlc.arg('to_time'), '[)')
ORDER BY lower(a.duration), a.id;

-- name: SearchAppointmentsByCustomerNameAndDealershipID :many
SELECT
  a.id,
  a.dealership_id,
  d.name AS dealership_name,
  d.open_time AS dealership_open_time,
  d.close_time AS dealership_close_time,
  a.service_id,
  sv.name AS service_name,
  sv.anticipated_minutes,
  a.bay_id,
  sb.name AS bay_name,
  a.technician_id,
  t.name AS technician_name,
  a.customer_name,
  a.status,
  a.duration,
  lower(a.duration) AS start_time,
  upper(a.duration) AS end_time,
  a.created_at,
  a.updated_at
FROM appointments a
JOIN dealerships d ON d.id = a.dealership_id
JOIN services sv ON sv.id = a.service_id
LEFT JOIN technicians t ON t.id = a.technician_id
LEFT JOIN service_bays sb ON sb.id = a.bay_id
WHERE a.dealership_id = sqlc.arg('dealership_id')
  AND unaccent(a.customer_name) ILIKE unaccent('%' || sqlc.arg('customer_name') || '%')
ORDER BY lower(a.duration) DESC, a.id DESC;

-- name: UpdateAppointmentStatusByID :one
UPDATE appointments
SET status = sqlc.arg('status')::status_type,
    updated_at = now()
WHERE id = sqlc.arg('appointment_id')
RETURNING id, status, updated_at;

-- name: MarkNoShowAppointmentsForDealershipInTimeRange :exec
-- End-of-day: mark appointments as no_show if they did not start/complete.
UPDATE appointments
SET status = 'no_show',
    updated_at = now()
WHERE id = ANY(sqlc.arg('appointment_ids')::int[])
  AND status NOT IN ('in_progress', 'completed', 'cancelled', 'no_show');

-- name: CreateAppointment :one
INSERT INTO appointments (
  dealership_id,
  service_id,
  bay_id,
  technician_id,
  customer_name,
  duration
)
VALUES (
  sqlc.arg('dealership_id'),
  sqlc.arg('service_id'),
  sqlc.arg('bay_id'),
  sqlc.arg('technician_id'),
  sqlc.arg('customer_name'),
  tstzrange(sqlc.arg('start_time')::timestamptz, sqlc.arg('end_time')::timestamptz, '[)')
)
RETURNING
  id,
  dealership_id,
  service_id,
  bay_id,
  technician_id,
  customer_name,
  status,
  duration,
  created_at,
  updated_at;

