
CREATE TABLE IF NOT EXISTS "delivery_service"."public"."delivery" (
   id SERIAL,
   tracking_code VARCHAR(50) NOT NULL,
   source_address TEXT NOT NULL,
   destination_address TEXT NOT NULL,
   status VARCHAR(50) NOT NULL,
--    status VARCHAR(255) CHECK
--        (status IN('CONFIRMED', 'IN_TRANSIT')) NOT NULL,
   created TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
   modified TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
   CONSTRAINT item PRIMARY KEY (id)
);
