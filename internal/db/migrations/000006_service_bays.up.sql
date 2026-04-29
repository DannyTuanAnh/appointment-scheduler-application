-- Create: service_bays
CREATE TABLE IF NOT EXISTS service_bays (
    id SERIAL PRIMARY KEY,
    dealership_id INT NOT NULL REFERENCES dealerships(id) ON DELETE CASCADE,
    bay_type_id INT NOT NULL REFERENCES service_bay_types(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ux_service_bays_dealership_name UNIQUE (dealership_id, name)
);

CREATE INDEX IF NOT EXISTS ix_service_bays_dealership_id ON service_bays (dealership_id);
CREATE INDEX IF NOT EXISTS ix_service_bays_type ON service_bays (bay_type_id);
CREATE INDEX IF NOT EXISTS ix_service_bays_active ON service_bays (dealership_id, is_active);
