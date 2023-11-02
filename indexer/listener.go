package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/NethermindEth/starknet.go/rpc"
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

		// reorg data prevention hard sync
		if ix.LastQueried.Timestamp.Add(time.Hour).Before(time.Now().UTC()) {
			logger.Info("hard sync at: " + fmt.Sprint(bInfo.BlockNumber))

			go ix.UpdateByFnsAll(bInfo.BlockNumber)

			if err := updateIxStatusDB(ix.LastQueried, *bInfo, ix.Store); err != nil {
				logger.Error(err, "cannot update the indexer status")
			}
		}

		if bInfo.BlockNumber > ix.LastQueried.BlockNumber {
			logger.Info("new block catched: " + fmt.Sprint(bInfo.BlockNumber))

			// Update the pools that not emit events
			go ix.UpdateByFns(bInfo.BlockNumber)

			// Update the pools whose events are caugh
			err := ix.GetEvents(ix.LastQueried.BlockNumber+1, *bInfo)
			if err != nil {
				logger.Error(err, "cannot get the events")
				continue
			}

			if err := updateIxStatusDB(ix.LastQueried, *bInfo, ix.Store); err != nil {
				logger.Error(err, "cannot update the indexer status")
			}
		}
	}
}

func updateIxStatusDB(lqs *status, bi rpc.BlockHashAndNumberOutput, store db.Store) error {
	lqs.BlockHash, lqs.BlockNumber = bi.BlockHash.String(), bi.BlockNumber
	lqs.Timestamp = time.Now().UTC()

	_, err := store.UpdateIndexerStatus(
		context.Background(), db.UpdateIndexerStatusParams{
			LastQueriedBlock: sql.NullInt64{Int64: int64(lqs.BlockNumber), Valid: true},
			LastQueriedHash:  sql.NullString{String: lqs.BlockHash, Valid: true},
		})

	return err
}
