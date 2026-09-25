-- ground-turn schema. The backend also runs GORM AutoMigrate on boot, which
-- adds missing columns/indexes idempotently; this file gives the MySQL 8
-- instance an explicit baseline on first `docker compose up`.

CREATE TABLE IF NOT EXISTS flight_turnaround (
  id INT AUTO_INCREMENT PRIMARY KEY,
  flight_no VARCHAR(32) NOT NULL,
  aircraft_reg VARCHAR(32) NULL,
  stand_no VARCHAR(32) NULL,
  arrival_time DATETIME(3) NULL,
  departure_time DATETIME(3) NULL,
  turnaround_status VARCHAR(32) NOT NULL,
  delay_reason VARCHAR(512) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  INDEX idx_flight_status (turnaround_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ground_resource (
  id INT AUTO_INCREMENT PRIMARY KEY,
  resource_code VARCHAR(64) NOT NULL,
  resource_type VARCHAR(32) NOT NULL,
  location VARCHAR(64) NULL,
  availability_status VARCHAR(32) NOT NULL,
  maintenance_due_at DATETIME(3) NULL,
  owner_team VARCHAR(64) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  UNIQUE INDEX uq_resource_code (resource_code),
  INDEX idx_resource_type (resource_type),
  INDEX idx_resource_status (availability_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ground_task (
  id INT AUTO_INCREMENT PRIMARY KEY,
  turnaround_id INT NOT NULL,
  task_type VARCHAR(32) NOT NULL,
  team_id INT NULL,
  planned_start DATETIME(3) NULL,
  deadline DATETIME(3) NULL,
  actual_finish DATETIME(3) NULL,
  status VARCHAR(32) NOT NULL,
  blocker_note VARCHAR(512) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  INDEX idx_task_turnaround (turnaround_id),
  INDEX idx_task_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS resource_booking (
  id INT AUTO_INCREMENT PRIMARY KEY,
  resource_id INT NOT NULL,
  turnaround_id INT NOT NULL,
  task_id INT NOT NULL,
  start_time DATETIME(3) NOT NULL,
  end_time DATETIME(3) NOT NULL,
  booking_status VARCHAR(32) NOT NULL,
  conflict_reason VARCHAR(512) NULL,
  -- Rebind history chain: old row -> RELEASED + replaced_by_id, new row
  -- carries replaced_booking_id back, so the before/after pair is
  -- queryable from either side.
  replaced_booking_id INT NULL,
  replaced_by_id INT NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  INDEX idx_booking_resource (resource_id),
  INDEX idx_booking_turnaround (turnaround_id),
  INDEX idx_booking_task (task_id),
  INDEX idx_booking_status (booking_status),
  INDEX idx_booking_start (start_time),
  INDEX idx_booking_replaced (replaced_booking_id),
  INDEX idx_booking_replaced_by (replaced_by_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS delay_event (
  id INT AUTO_INCREMENT PRIMARY KEY,
  turnaround_id INT NOT NULL,
  delay_type VARCHAR(32) NULL,
  minutes INT NULL,
  root_cause VARCHAR(512) NULL,
  responsibility_team VARCHAR(64) NULL,
  resolved_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  INDEX idx_delay_turnaround (turnaround_id),
  INDEX idx_delay_resolved (resolved_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_log (
  id INT AUTO_INCREMENT PRIMARY KEY,
  actor VARCHAR(64) NOT NULL,
  action VARCHAR(64) NOT NULL,
  target_type VARCHAR(64) NOT NULL,
  target_id VARCHAR(64) NULL,
  detail VARCHAR(1024) NULL,
  created_at DATETIME(3) NULL,
  INDEX idx_audit_action (action),
  INDEX idx_audit_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
