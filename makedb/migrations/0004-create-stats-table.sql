CREATE TABLE IF NOT EXISTS button_event (
    id SERIAL PRIMARY KEY,
    x_coord INTEGER NOT NULL,
    y_coord INTEGER NOT NULL,
    button_id INTEGER NOT NULL,
    event_type VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE INDEX IF NOT EXISTS idx_button_event_event_type 
    ON button_event (event_type, created_at); 

CREATE TABLE IF NOT EXISTS button_stat (
    stat_key varchar(32) NOT NULL PRIMARY KEY,
    stat_name varchar(255) NOT NULL,
    stat_desc varchar(2000) NOT NULL,
    val INTEGER NOT NULL DEFAULT 0,
    "scale" INTEGER NOT NULL DEFAULT 1,
    "order" INTEGER NOT NULL
);

INSERT INTO button_stat (stat_key, stat_name, stat_desc, val, "scale", "order") 
VALUES 
('buttons_pressed', 'Buttons Pressed', 'Total number of buttons that users have pressed', 0, 1, 1),
('pressed_last_day', 'Buttons pressed in the last day', 'Total number of buttons that users have pressed in the last 24 hours', 0, 1, 2),
('presses_per_second', 'Presses per second', 'Average number of button presses per second', 0, -3, 3)
ON CONFLICT (stat_key) DO NOTHING;

DO $$
BEGIN
DROP PROCEDURE IF EXISTS update_button_stats;

CREATE OR REPLACE PROCEDURE update_button_stats()
AS $BODY$
BEGIN

UPDATE button_stat SET val = (
    SELECT COUNT(*) 
    FROM button_event
    WHERE event_type = 'press'
) 
WHERE stat_key = 'buttons_pressed';


UPDATE button_stat SET val = (
    SELECT COUNT(*) 
    FROM button_event
    WHERE event_type = 'press'
    AND created_at >= NOW() - INTERVAL '1 day'
) 
WHERE stat_key = 'pressed_last_day';


UPDATE button_stat SET val = (
    SELECT COUNT(*) * 1000 / 600
    FROM button_event
    WHERE event_type = 'press'
    AND created_at >= NOW() - INTERVAL '10 minutes'
) 
WHERE stat_key = 'presses_per_second';

END;
$BODY$ LANGUAGE PLPGSQL;

END $$;