CREATE TYPE technician_level AS ENUM ('fresher', 'junior', 'senior');

-- Create: technicians
CREATE TABLE IF NOT EXISTS technicians (
    id SERIAL PRIMARY KEY,
    dealership_id INT NOT NULL REFERENCES dealerships(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    level technician_level NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    inactive_since TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_technicians_dealership_id ON technicians (dealership_id);

