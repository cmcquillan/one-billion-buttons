DO $$
BEGIN

DROP PROCEDURE IF EXISTS public.set_button_color;

DROP TABLE IF EXISTS public.button;

DROP DOMAIN IF EXISTS fixed_bytea;

END $$;
