-- name: CreateToken :one
INSERT INTO tokens (
  address,
  name,
  symbol,
  decimals,
  ticker
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateBaseNativeStatus :one
UPDATE tokens
SET base = $2, native = $3
WHERE address = $1
RETURNING *;

-- name: UpdateTicker :one
UPDATE tokens
SET ticker = $2
WHERE address = $1
RETURNING *;

-- name: GetTokenByAddress :one
SELECT * FROM tokens
WHERE address = $1 LIMIT 1;

-- name: GetTokenBySymbol :one
SELECT * FROM tokens
WHERE symbol = $1 LIMIT 1;

-- name: GetBaseTokens :many
SELECT * FROM tokens
WHERE base = true
ORDER BY name;

-- name: GetNativeTokens :many
SELECT * FROM tokens
WHERE native = true
ORDER BY name;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE address = $1;