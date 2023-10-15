-- name: CreateAmm :one
INSERT INTO amms (
  dex_name,
  router_address,
  key,
  algorithm_type
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAmmById :one
SELECT * FROM amms
WHERE amm_id = $1 LIMIT 1;

-- name: GetAllAmms :many
SELECT * FROM amms
ORDER BY amm_id;

-- name: GetKeys :many
SELECT DISTINCT key FROM amms
ORDER BY key;

-- name: GetAmmByDEX :many
SELECT * FROM amms
WHERE dex_name = $1
ORDER BY dex_name;

-- name: GetAmmKeys :many
SELECT DISTINCT key FROM amms
WHERE key IS NOT NULL AND key != ''
ORDER BY KEY;

-- name: DeleteAmm :exec
DELETE FROM amms
WHERE amm_id = $1;