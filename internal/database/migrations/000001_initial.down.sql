ALTER TABLE vms DROP CONSTRAINT IF EXISTS fk_vms_task_id;

DROP INDEX IF EXISTS idx_vms_name_active;
DROP INDEX IF EXISTS idx_vms_deleted_at;
DROP INDEX IF EXISTS idx_vms_task_id;

DROP INDEX IF EXISTS idx_tasks_vm_id_created_at;
DROP INDEX IF EXISTS idx_tasks_created_at;

DROP TABLE IF EXISTS tasks CASCADE;

DROP TABLE IF EXISTS vms CASCADE;