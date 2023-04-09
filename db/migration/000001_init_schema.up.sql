-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-04-09T11:37:47.631Z

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

CREATE TABLE "pools_v2" (
  "address" varchar PRIMARY KEY NOT NULL,
  "amm_id" bigint NOT NULL,
  "token_a" varchar NOT NULL,
  "token_b" varchar NOT NULL,
  "reserve_a" numeric NOT NULL DEFAULT 0,
  "reserve_b" numeric NOT NULL DEFAULT 0,
  "fee" numeric NOT NULL,
  "total_value" numeric NOT NULL DEFAULT 0,
  "last_updated" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "pools_v3" (
  "address" varchar PRIMARY KEY NOT NULL,
  "amm_id" bigint NOT NULL,
  "token_a" varchar NOT NULL,
  "token_b" varchar NOT NULL,
  "sqrt_price_X96" varchar NOT NULL,
  "tick" varchar NOT NULL,
  "observation_index" varchar NOT NULL,
  "observation_cardinality" varchar NOT NULL,
  "observation_cardinality_next" varchar NOT NULL,
  "fee" numeric NOT NULL,
  "unlocked" bool NOT NULL,
  "liquidity_gross" varchar NOT NULL DEFAULT '0',
  "liquidity_net" varchar NOT NULL DEFAULT '0',
  "fee_growth_outside_0X128" varchar NOT NULL DEFAULT '0',
  "fee_growth_outside_1X128" varchar NOT NULL DEFAULT '0',
  "tick_cumulative_outside" varchar NOT NULL DEFAULT '0',
  "seconds_per_liquidity_outside_X128" varchar NOT NULL DEFAULT '0',
  "seconds_outside" varchar NOT NULL DEFAULT '0',
  "initialized" bool NOT NULL DEFAULT false,
  "liquidity" varchar NOT NULL
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
  "last_queried" bigint DEFAULT 0,
  "last_updated" timestamptz DEFAULT '0001-01-01'
);

COMMENT ON COLUMN "pools_v3"."sqrt_price_X96" IS 'sqrtPriceX96';

COMMENT ON COLUMN "pools_v3"."tick" IS 'tick';

COMMENT ON COLUMN "pools_v3"."observation_index" IS 'observationIndex';

COMMENT ON COLUMN "pools_v3"."observation_cardinality" IS 'observationCardinality';

COMMENT ON COLUMN "pools_v3"."observation_cardinality_next" IS 'observationCardinalityNext';

COMMENT ON COLUMN "pools_v3"."fee" IS 'feeProtocol';

COMMENT ON COLUMN "pools_v3"."unlocked" IS 'unlocked';

COMMENT ON COLUMN "pools_v3"."liquidity_gross" IS 'liquidityGross';

COMMENT ON COLUMN "pools_v3"."liquidity_net" IS 'liquidityNet';

COMMENT ON COLUMN "pools_v3"."fee_growth_outside_0X128" IS 'feeGrowthOutside0X128';

COMMENT ON COLUMN "pools_v3"."fee_growth_outside_1X128" IS 'feeGrowthOutside1X128';

COMMENT ON COLUMN "pools_v3"."tick_cumulative_outside" IS 'tickCumulativeOutside';

COMMENT ON COLUMN "pools_v3"."seconds_per_liquidity_outside_X128" IS 'secondsPerLiquidityOutsideX128';

COMMENT ON COLUMN "pools_v3"."seconds_outside" IS 'secondsOutside';

COMMENT ON COLUMN "pools_v3"."initialized" IS 'initialized';

COMMENT ON COLUMN "pools_v3"."liquidity" IS 'liquidity';

COMMENT ON COLUMN "amms"."created_at" IS 'initialized';

ALTER TABLE "pools_v2" ADD FOREIGN KEY ("amm_id") REFERENCES "amms" ("amm_id");

ALTER TABLE "pools_v2" ADD FOREIGN KEY ("token_a") REFERENCES "tokens" ("address");

ALTER TABLE "pools_v2" ADD FOREIGN KEY ("token_b") REFERENCES "tokens" ("address");

ALTER TABLE "pools_v3" ADD FOREIGN KEY ("amm_id") REFERENCES "amms" ("amm_id");

ALTER TABLE "pools_v3" ADD FOREIGN KEY ("token_a") REFERENCES "tokens" ("address");

ALTER TABLE "pools_v3" ADD FOREIGN KEY ("token_b") REFERENCES "tokens" ("address");
