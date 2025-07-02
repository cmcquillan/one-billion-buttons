DO $$
BEGIN
DROP PROCEDURE IF EXISTS set_button_color;

/*
 * Example: call set_button_color (1, 1, 56, '\xFFAB03');
 * select buttons from button where x_coord = 1 and y_coord = 1;
 */
CREATE OR REPLACE PROCEDURE set_button_color(
    x INTEGER,
    y INTEGER,
    ix INTEGER,
    rgbVal BYTEA)
AS $BODY$
DECLARE
    ixs INTEGER;
BEGIN

ixs := ix * 3;

UPDATE button SET 
    buttons = set_byte(set_byte(set_byte(buttons, ixs + 2, get_byte(rgbVal, 2)), ixs + 1, get_byte(rgbVal, 1)), ixs, get_byte(rgbVal, 0))
    ,version = version + 1
WHERE 
    x_coord = x AND
    y_coord = y AND
    substring(buttons FROM (ixs+1) FOR 3) = '\x000000';

END;
$BODY$ LANGUAGE PLPGSQL;

END $$;