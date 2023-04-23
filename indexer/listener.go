package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"

	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *indexer) ListenBlocks() {
	// Connect to socket server
	conn, err := net.Dial("tcp", ix.config.SocketAddress)
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
			continue
		}

		bn, err := strconv.Atoi(string(buffer[:n]))
		if err != nil {
			logger.Error(err, "cannot convert event to block number")
		}

		ubn := uint64(bn)

		// FIXME: temporary solution for the late sync. problem in the issue #14
		time.Sleep(time.Second)

		if ubn > *ix.lastQueried {
			logger.Info("new block catched: " + fmt.Sprint(bn))

			err := ix.getEvents(*ix.lastQueried+1, ubn)
			if err != nil {
				logger.Error(err, "cannot get the events")
				return
			}

			ix.lastQueried = &ubn
			_, err = ix.store.UpdateIndexerStatus(context.Background(), sql.NullInt64{Int64: int64(*ix.lastQueried), Valid: true})
			if err != nil {
				logger.Error(err, "cannot update the indexer status")
			}
		}
	}
}
