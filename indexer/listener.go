package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"time"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *Indexer) ListenBlocks() {
	// Connect to socket server
	conn, err := net.Dial("tcp", ix.Config.SocketAddress)
	if err != nil {
		logger.Error(err, "cannot connect to the socket server")

		time.Sleep(3 * time.Second)
		go ix.ListenBlocks()
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Error(err, "cannot read the socket server")

			time.Sleep(3 * time.Second)
			go ix.ListenBlocks()
			return
		}

		var bInfo *rpc.BlockHashAndNumberOutput = &rpc.BlockHashAndNumberOutput{}
		if err := json.Unmarshal(buffer[:n], bInfo); err != nil {
			logger.Error(err, "cannot convert event to block number")
		}

		// FIXME: temporary solution for the late sync. problem in the issue #14
		time.Sleep(time.Second)

		if bInfo.BlockNumber > ix.LastQueried.BlockNumber {
			logger.Info("new block catched: " + fmt.Sprint(bInfo.BlockNumber))

			err := ix.GetEvents(ix.LastQueried.BlockNumber+1, bInfo.BlockNumber)
			if err != nil {
				logger.Error(err, "cannot get the events")
				return
			}

			ix.LastQueried = bInfo
			_, err = ix.Store.UpdateIndexerStatus(
				context.Background(), db.UpdateIndexerStatusParams{
					LastQueriedBlock: sql.NullInt64{Int64: int64(ix.LastQueried.BlockNumber), Valid: true},
					LastQueriedHash: sql.NullString{String: ix.LastQueried.BlockHash, Valid: true},
				})
			if err != nil {
				logger.Error(err, "cannot update the indexer status")
			}
		}
	}
}
