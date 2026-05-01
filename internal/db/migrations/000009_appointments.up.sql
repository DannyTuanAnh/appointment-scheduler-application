-- Requires for GiST indexing on scalar types in multi-column GiST and for exclusion constraints
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Appointment status enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
        CREATE TYPE status_type AS ENUM ('confirmed', 'in_progress', 'completed', 'cancelled', 'no_show');
    END IF;
END$$;

-- Create: appointments
CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    dealership_id INT NOT NULL REFERENCES dealerships(id) ON DELETE CASCADE,
    service_id INT NOT NULL REFERENCES services(id) ON DELETE SET NULL,

    -- Assigned resources (can be nullable until scheduling occurs)
    bay_id INT NOT NULL REFERENCES service_bays(id) ON DELETE SET NULL,
    technician_id INT NOT NULL REFERENCES technicians(id) ON DELETE SET NULL,

    customer_name VARCHAR(100) NOT NULL DEFAULT 'confirmed',

    status status_type NOT NULL,

    -- Time window: inclusive start, exclusive end: [start, end)
    duration TSTZRANGE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Ensure valid range and convenient constraints for planners
    CONSTRAINT ck_appointments_duration_valid CHECK (
        NOT isempty(duration)
        AND lower(duration) IS NOT NULL
        AND upper(duration) IS NOT NULL
        AND lower(duration) < upper(duration)
        AND lower_inc(duration) --- start at with [ that including that time in the appointment 
        AND NOT upper_inc(duration) --- end with ) that excluding that time in the appointment
    )
);

-- Helpful indexes for time-window queries
CREATE INDEX IF NOT EXISTS ix_appointments_dealership_duration_gist
    ON appointments USING GIST (dealership_id, duration);

CREATE INDEX IF NOT EXISTS ix_appointments_status_duration_gist
    ON appointments USING GIST (status, duration);

-- Partial GiST indexes for availability lookups
CREATE INDEX IF NOT EXISTS ix_appointments_bay_duration_gist
    ON appointments USING GIST (bay_id, duration)
    WHERE bay_id IS NOT NULL AND status <> 'cancelled';

CREATE INDEX IF NOT EXISTS ix_appointments_technician_duration_gist
    ON appointments USING GIST (technician_id, duration)
    WHERE technician_id IS NOT NULL AND status <> 'cancelled';

-- Prevent overlapping bookings at DB layer (ignore cancelled appointments)
-- Bay overlap prevention
ALTER TABLE appointments
    ADD CONSTRAINT ex_appointments_no_overlap_bay
    EXCLUDE USING GIST (
        bay_id WITH =,
        duration WITH &&
    )
    WHERE (bay_id IS NOT NULL AND status <> 'cancelled');

-- Technician overlap prevention
ALTER TABLE appointments
    ADD CONSTRAINT ex_appointments_no_overlap_technician
    EXCLUDE USING GIST (
        technician_id WITH =,
        duration WITH &&
    )
    WHERE (technician_id IS NOT NULL AND status <> 'cancelled');
