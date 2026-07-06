CREATE TABLE IF NOT EXISTS vms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    cpu INTEGER NOT NULL,
    ram_mb INTEGER NOT NULL,
    disk_gb INTEGER NOT NULL,
    task_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_vms_status CHECK (status IN ('pending', 'running', 'stopped')),
    CONSTRAINT chk_vms_cpu CHECK (cpu > 0),
    CONSTRAINT chk_vms_ram_mb CHECK (ram_mb > 0),
    CONSTRAINT chk_vms_disk_gb CHECK (disk_gb > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_vms_name_active ON vms (name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_vms_deleted_at ON vms (deleted_at);
CREATE INDEX IF NOT EXISTS idx_vms_task_id ON vms (task_id);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    vm_id BIGINT NOT NULL REFERENCES vms(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_tasks_type CHECK (type IN ('provision', 'start', 'stop', 'delete')),
    CONSTRAINT chk_tasks_status CHECK (status IN ('pending', 'running', 'completed', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_tasks_vm_id_created_at ON tasks (vm_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks (created_at DESC);

ALTER TABLE vms
    ADD CONSTRAINT fk_vms_task_id
    FOREIGN KEY (task_id) REFERENCES tasks(id)
    ON DELETE SET NULL
    DEFERRABLE INITIALLY DEFERRED;
