CREATE TABLE IF NOT EXISTS components (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    type text NOT NULL,
    manufacturer text NOT NULL,
    model text NOT NULL,
    price numeric(10, 2) NOT NULL,
    rating smallint NOT NULL DEFAULT 0,
    image_path text NOT NULL DEFAULT '/static/images/placeholder.webp' 
);