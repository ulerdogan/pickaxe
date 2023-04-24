package socket

import (
	"strconv"

	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func setupJobs(sc *socket) {
	sc.scheduler.Every(5).Seconds().Do(sc.QueryBlocks)
}

func (sc *socket) QueryBlocks() {
	sc.scMutex.Lock()
	defer sc.scMutex.Unlock()

	block, err := sc.client.LastBlock()
	if err != nil {
		logger.Error(err, "cannot get the last block")
		return
	}

	bn := block.BlockNumber

	if sc.blockInfo == nil || bn > sc.blockInfo.BlockNumber {
		sc.blockInfo = block
		logger.Info("new block cathed: " + strconv.Itoa(int(bn)))
	}
}
