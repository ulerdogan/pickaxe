-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-09-18T01:11:09.929Z

CREATE TABLE "tokens" (
  "address" varchar PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "symbol" varchar NOT NULL,
  "decimals" int NOT NULL,
  "base" bool NOT NULL DEFAULT false,
  "native" bool NOT NULL DEFAULT false,
  "ticker" varchar NOT NULL DEFAULT '',
  "price" varchar NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "pools" (
  "pool_id" bigserial PRIMARY KEY,
  "address" varchar NOT NULL,
  "amm_id" bigint NOT NULL,
  "token_a" varchar NOT NULL,
  "token_b" varchar NOT NULL,
  "reserve_a" numeric NOT NULL DEFAULT 0,
  "reserve_b" numeric NOT NULL DEFAULT 0,
  "fee" varchar NOT NULL DEFAULT '',
  "total_value" numeric NOT NULL DEFAULT 0,
  "extra_data" varchar,
  "general_extra_data" varchar,
  "last_updated" timestamptz NOT NULL DEFAULT '0001-01-01',
  "last_block" bigint NOT NULL DEFAULT 0
);

CREATE TABLE "amms" (
  "amm_id" bigserial PRIMARY KEY,
  "dex_name" varchar NOT NULL,
  "router_address" varchar NOT NULL,
  "key" varchar NOT NULL DEFAULT '',
  "algorithm_type" varchar NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "indexer" (
  "id" int PRIMARY KEY NOT NULL,
  "hashed_password" varchar NOT NULL,
  "last_queried_block" bigint DEFAULT 0,
  "last_queried_hash" varchar DEFAULT '',
  "last_updated" timestamptz DEFAULT '0001-01-01'
);

COMMENT ON COLUMN "amms"."created_at" IS 'initialized';

ALTER TABLE "pools" ADD FOREIGN KEY ("amm_id") REFERENCES "amms" ("amm_id");

ALTER TABLE "pools" ADD FOREIGN KEY ("token_a") REFERENCES "tokens" ("address");

ALTER TABLE "pools" ADD FOREIGN KEY ("token_b") REFERENCES "tokens" ("address");
