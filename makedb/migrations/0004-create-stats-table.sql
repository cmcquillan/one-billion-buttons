CREATE TABLE IF NOT EXISTS button_event (
    id SERIAL PRIMARY KEY,
    x_coord INTEGER NOT NULL,
    y_coord INTEGER NOT NULL,
    button_id INTEGER NOT NULL,
    event_type INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE IF NOT EXISTS button_stat (
    stat_key varchar(32) NOT NULL PRIMARY KEY,
    stat_name varchar(255) NOT NULL,
    stat_desc varchar(2000) NOT NULL,
    val INTEGER NOT NULL DEFAULT 0,
    "scale" INTEGER NOT NULL DEFAULT 1
);

INSERT INTO button_stat (stat_key, stat_name, stat_desc, val, "scale") 
VALUES 
('buttons_pressed', 'Buttons Pressed', 'Total number of buttons that users have pressed', 0, 1),
('pressed_last_day', 'Buttons pressed in the last day', 'Total number of buttons that users have pressed in the last 24 hours', 0, 1),
('presses_per_second', 'Presses per second', 'Average number of button presses per second', 0, 1)
ON CONFLICT (stat_key) DO NOTHING;