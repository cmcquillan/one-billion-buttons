DO $$
BEGIN

DROP PROCEDURE IF EXISTS set_button_color;
DROP FUNCTION IF EXISTS get_minimap_color;

CREATE OR REPLACE FUNCTION get_minimap_color(bytes fixed_bytea) 
RETURNS bytea
AS $BODY$
DECLARE 
    cr INTEGER := 0;
    cg INTEGER := 0;
    cb INTEGER := 0;
BEGIN

FOR i IN 0..297 BY 3 LOOP
    cr = cr + get_byte(bytes, i);
    cg = cg + get_byte(bytes, i+1);
    cb = cb + get_byte(bytes, i+2);
END LOOP;

    cr = cr / 100;
    cg = cg / 100;
    cb = cb / 100;

    RETURN to_hex(cr) || to_hex(cg) || to_hex(cb);
END;
$BODY$ LANGUAGE PLPGSQL;

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
    rV INTEGER := get_byte(rgbVal, 0);
    gV INTEGER := get_byte(rgbVal, 1);
    bV INTEGER := get_byte(rgbVal, 2);
BEGIN

ixs := ix * 3;

UPDATE button SET 
    buttons = set_byte(set_byte(set_byte(buttons, ixs + 2, bV), ixs + 1, gV), ixs, rV)
    ,version = version + 1
    ,map_value = get_minimap_color(buttons)
WHERE 
    x_coord = x AND
    y_coord = y AND
    substring(buttons FROM (ixs+1) FOR 3) = '\x000000';

END;
$BODY$ LANGUAGE PLPGSQL;

END $$;
