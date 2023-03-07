-- name: CreatePool :one
INSERT INTO pools (
  address,
  amm_id,
  token_a,
  token_b,
  reserve_a,
  reserve_b
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPoolByAddress :one
SELECT * FROM pools
WHERE address = $1 LIMIT 1;

-- name: GetPoolsByPair :many
SELECT * FROM pools
WHERE token_a = $1 AND token_b = $2
ORDER BY amm_id;

-- name: GetPoolsByToken :many
SELECT * FROM pools
WHERE token_a = $1 OR token_b = $1
ORDER BY amm_id;

-- name: GetPoolsByAmm :many
SELECT * FROM pools
WHERE amm_id = $1
ORDER BY address;

-- name: DeletePool :exec
DELETE FROM pools
WHERE address = $1;