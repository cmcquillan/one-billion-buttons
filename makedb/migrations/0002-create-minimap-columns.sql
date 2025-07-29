DO $$
BEGIN

ALTER TABLE button
    ADD COLUMN IF NOT EXISTS map_value bytea NOT NULL DEFAULT '\x000000'::bytea;

END;
$$;
