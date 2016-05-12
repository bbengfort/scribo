/**
 * 0001-initialize.sql
 * Copyright 2016 University of Maryland
 *
 * Author:  Benjamin Bengfort <benjamin@bengfort.com>
 * Created: Thu May 12 11:26:40 2016 -0400
 */

-------------------------------------------------------------------------
-- Ensure transaction security by placing all CREATE and ALTER statements
-- inside of `BEGIN` and `COMMIT` statements.
-------------------------------------------------------------------------

BEGIN;

/**
 *  CREATE ENTITY TABLES
 */

-------------------------------------------------------------------------
-- nodes Table
-------------------------------------------------------------------------

-- DROP TABLE IF EXISTS "nodes";

CREATE TABLE "nodes"
(
    "id" SERIAL NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "address" VARCHAR(45),
    "dns" VARCHAR(255),
    "key" VARCHAR(44),
    "created" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-------------------------------------------------------------------------
-- pings Table
-------------------------------------------------------------------------

-- DROP TABLE IF EXISTS "pings";

CREATE TABLE "pings"
(
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "source_id" INT NOT NULL,
    "target_id" INT NOT NULL,
    "payload" INT,
    "latency" DOUBLE PRECISION,
    "timeout" BOOLEAN,
    "created" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


/*
 *  ALTER TABLE ADD FOREIGN KEYS AFTER ENTITY TABLES
 */

 -------------------------------------------------------------------------
 -- pings.source_id -> node.id
 -------------------------------------------------------------------------

 ALTER TABLE "pings" ADD CONSTRAINT "fk_pings_source_id"
     FOREIGN KEY ("source_id")
     REFERENCES "nodes" ("id") MATCH SIMPLE
     DEFERRABLE INITIALLY DEFERRED;

     ---------------------------------------------------------------------
     -- pings.source_id Foreign Key Index
     ---------------------------------------------------------------------

     -- DROP INDEX IF EXISTS "idx_pings_source_id";

     CREATE INDEX "idx_pings_source_id"
         ON "pings" USING BTREE ("source_id");

 -------------------------------------------------------------------------
 -- pings.target_id -> node.id
 -------------------------------------------------------------------------

 ALTER TABLE "pings" ADD CONSTRAINT "fk_pings_target_id"
     FOREIGN KEY ("target_id")
     REFERENCES "nodes" ("id") MATCH SIMPLE
     DEFERRABLE INITIALLY DEFERRED;

     ---------------------------------------------------------------------
     -- pings.target_id Foreign Key Index
     ---------------------------------------------------------------------

     -- DROP INDEX IF EXISTS "idx_pings_target_id";

     CREATE INDEX "idx_pings_target_id"
         ON "pings" USING BTREE ("target_id");

 /**
  *  CREATE INDICIES
  */

 COMMIT;

 -------------------------------------------------------------------------
 -- No CREATE or ALTER statements should be outside of the `COMMIT`.
 -------------------------------------------------------------------------
