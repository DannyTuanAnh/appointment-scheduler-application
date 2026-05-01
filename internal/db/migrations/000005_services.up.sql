-- Create: services
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    required_bay_type_id INT NOT NULL REFERENCES service_bay_types(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    anticipated_minutes INT NOT NULL CHECK (anticipated_minutes > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_services_name ON services (name);
