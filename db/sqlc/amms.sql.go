// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: amms.sql

package db

import (
	"context"
)

const createAmm = `-- name: CreateAmm :one
INSERT INTO amms (
  dex_name,
  fee,
  router_address,
  algorithm_type
) VALUES (
  $1, $2, $3, $4
) RETURNING amm_id, dex_name, fee, router_address, algorithm_type
`

type CreateAmmParams struct {
	DexName       string `json:"dex_name"`
	Fee           string `json:"fee"`
	RouterAddress string `json:"router_address"`
	AlgorithmType string `json:"algorithm_type"`
}

func (q *Queries) CreateAmm(ctx context.Context, arg CreateAmmParams) (Amm, error) {
	row := q.db.QueryRowContext(ctx, createAmm,
		arg.DexName,
		arg.Fee,
		arg.RouterAddress,
		arg.AlgorithmType,
	)
	var i Amm
	err := row.Scan(
		&i.AmmID,
		&i.DexName,
		&i.Fee,
		&i.RouterAddress,
		&i.AlgorithmType,
	)
	return i, err
}

const deleteAmm = `-- name: DeleteAmm :exec
DELETE FROM amms
WHERE amm_id = $1
`

func (q *Queries) DeleteAmm(ctx context.Context, ammID int64) error {
	_, err := q.db.ExecContext(ctx, deleteAmm, ammID)
	return err
}

const getAmmByDEX = `-- name: GetAmmByDEX :many
SELECT amm_id, dex_name, fee, router_address, algorithm_type FROM amms
WHERE dex_name = $1
ORDER BY dex_name
`

func (q *Queries) GetAmmByDEX(ctx context.Context, dexName string) ([]Amm, error) {
	rows, err := q.db.QueryContext(ctx, getAmmByDEX, dexName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Amm{}
	for rows.Next() {
		var i Amm
		if err := rows.Scan(
			&i.AmmID,
			&i.DexName,
			&i.Fee,
			&i.RouterAddress,
			&i.AlgorithmType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAmmById = `-- name: GetAmmById :one
SELECT amm_id, dex_name, fee, router_address, algorithm_type FROM amms
WHERE amm_id = $1 LIMIT 1
`

func (q *Queries) GetAmmById(ctx context.Context, ammID int64) (Amm, error) {
	row := q.db.QueryRowContext(ctx, getAmmById, ammID)
	var i Amm
	err := row.Scan(
		&i.AmmID,
		&i.DexName,
		&i.Fee,
		&i.RouterAddress,
		&i.AlgorithmType,
	)
	return i, err
}