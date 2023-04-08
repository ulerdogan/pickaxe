-- name: InitIndexer :one
INSERT INTO indexer (
  id,
  last_queried,
  last_updated
) VALUES (
  0, $1, NOW()
) RETURNING *;

-- name: GetIndexerStatus :one
SELECT * FROM indexer
WHERE id = 0 LIMIT 1;

-- name: UpdateIndexerStatus :one
UPDATE indexer
SET last_queried = $1 AND last_updated = NOW()
WHERE id = 0
RETURNING *;