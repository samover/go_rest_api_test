-- Migration: products
-- Created at: 2017-05-27 11:52:15
-- ====  UP  ====

BEGIN;
    CREATE TABLE IF NOT EXISTS products
    (
        id SERIAL,
        name TEXT NOT NULL,
        price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
        CONSTRAINT products_pkey PRIMARY KEY (id)
    );

COMMIT;

-- ==== DOWN ====

BEGIN;
    DROP TABLE products

COMMIT;
