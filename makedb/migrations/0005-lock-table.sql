DO $$
BEGIN

CREATE TABLE IF NOT EXISTS sync_lock (
    id varchar(32) PRIMARY KEY,
    lock_val UUID NULL,
    lock_time TIMESTAMP NULL DEFAULT NULL
);

END $$;
