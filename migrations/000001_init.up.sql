CREATE TABLE IF NOT EXISTS public.goods (
    code serial PRIMARY KEY,
    name text,
    size numeric,
    value integer,
    stock_id uuid
);

CREATE TABLE IF NOT EXISTS public.stocks (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" text NULL,
	available bool NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS public.res_cen (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	good_code int4 NULL,
	stock_id uuid NOT NULL,
	value int4 NULL
);