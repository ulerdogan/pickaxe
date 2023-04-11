-- name: CreatePool :one
INSERT INTO pools_v2 (
  address,
  amm_id,
  token_a,
  token_b,
  fee
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPoolByAddress :one
SELECT * FROM pools_v2
WHERE address = $1 LIMIT 1;

-- name: GetPoolsByPair :many
SELECT * FROM pools_v2
WHERE token_a = $1 AND token_b = $2
ORDER BY amm_id;

-- name: GetPoolsByToken :many
SELECT * FROM pools_v2
WHERE token_a = $1 OR token_b = $1
ORDER BY amm_id;

-- name: GetPoolsByAmm :many
SELECT * FROM pools_v2
WHERE amm_id = $1
ORDER BY address;

-- name: DeletePool :exec
DELETE FROM pools_v2
WHERE address = $1;