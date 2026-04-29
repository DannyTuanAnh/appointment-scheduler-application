-- Create: dealerships
CREATE TABLE IF NOT EXISTS dealerships (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    open_time TIME NOT NULL, 
    close_time TIME NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_dealerships_name ON dealerships (name);
