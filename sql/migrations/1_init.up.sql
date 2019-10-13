CREATE TABLE IF NOT EXISTS pr_events (
  `pr_number`  INT,
  `repository` VARCHAR(255),
  PRIMARY KEY (`pr_number`, `repository`)
);
