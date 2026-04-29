-- Create: technician_skills
CREATE TABLE IF NOT EXISTS technician_skills (
    technician_id INT NOT NULL REFERENCES technicians(id) ON DELETE CASCADE,
    skill_id INT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_technician_skills PRIMARY KEY (technician_id, skill_id)
);

CREATE INDEX IF NOT EXISTS ix_technician_skills_technician_id ON technician_skills (technician_id);
CREATE INDEX IF NOT EXISTS ix_technician_skills_skill_id ON technician_skills (skill_id);
