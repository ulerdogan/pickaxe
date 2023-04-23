package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
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

		bn, err := strconv.Atoi(string(buffer[:n]))
		if err != nil {
			logger.Error(err, "cannot convert event to block number")
		}

		ubn := uint64(bn)

		// FIXME: temporary solution for the late sync. problem in the issue #14
		time.Sleep(time.Second)

		if ubn > ix.LastQueried.BlockNumber {
			logger.Info("new block catched: " + fmt.Sprint(bn))

			err := ix.GetEvents(ix.LastQueried.BlockNumber+1, ubn)
			if err != nil {
				logger.Error(err, "cannot get the events")
				return
			}

			ix.LastQueried = &rpc.BlockHashAndNumberOutput{BlockNumber: ubn}
			_, err = ix.Store.UpdateIndexerStatus(context.Background(), sql.NullInt64{Int64: int64(ix.LastQueried.BlockNumber), Valid: true})
			if err != nil {
				logger.Error(err, "cannot update the indexer status")
			}
		}
	}
}
