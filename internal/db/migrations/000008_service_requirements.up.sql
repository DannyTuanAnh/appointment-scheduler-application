-- Create: service_requirements
CREATE TABLE IF NOT EXISTS service_requirements (
    service_id INT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    skill_id INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_service_requirements PRIMARY KEY (service_id, skill_id)
);

CREATE INDEX IF NOT EXISTS ix_service_requirements_service_id ON service_requirements (service_id);
CREATE INDEX IF NOT EXISTS ix_service_requirements_skill_id ON service_requirements (skill_id);
