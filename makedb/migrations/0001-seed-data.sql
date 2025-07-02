DO $$
DECLARE 
    rowCount INT := (SELECT COUNT(*) FROM button);
BEGIN

IF rowCount = 0 THEN
    FOR x IN 1..4096 LOOP
        INSERT INTO button (x_coord, y_coord)
        SELECT x, n.y
        FROM generate_series(1, 2500) AS n(y);
    END LOOP;
END IF;

END;
$$;