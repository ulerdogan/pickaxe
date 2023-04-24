-- name: InitIndexer :one
INSERT INTO indexer (
  id,
  hashed_password,
  last_queried_block,
  last_queried_hash,
  last_updated
) VALUES (
  0, $1, $2, $3, NOW()
) RETURNING *;

-- name: GetIndexerStatus :one
SELECT * FROM indexer
WHERE id = 0 LIMIT 1;

-- name: GetHashedIndexerPwd :one
SELECT hashed_password FROM indexer
WHERE id = 0 LIMIT 1;

-- name: UpdateIndexerStatus :one
UPDATE indexer
SET last_queried_block = $1, last_queried_hash = $2, last_updated = NOW()
WHERE id = 0
RETURNING *;