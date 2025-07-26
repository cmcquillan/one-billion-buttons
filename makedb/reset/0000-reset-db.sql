DO $$
BEGIN

DROP PROCEDURE IF EXISTS public.set_button_color;

DROP PROCEDURE IF EXISTS public.update_button_stats;

DROP TABLE IF EXISTS public.button_event;

DROP TABLE IF EXISTS public.button_stat;

DROP TABLE IF EXISTS public.button;

DROP DOMAIN IF EXISTS fixed_bytea;

END $$;
