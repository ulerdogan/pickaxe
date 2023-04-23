package socket

import (
	"encoding/json"
	"time"

	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (sc *socket) Sync() {
	conn, err := sc.listener.Accept()
	if err != nil {
		logger.Error(err, "error accepted in the listener")
		return
	} else {
		go sc.Sync()
	}
	defer conn.Close()
	var lastSent uint64 = 0

	for {
		if lastSent == 0 || sc.blockInfo.BlockNumber > lastSent {
			lastSent = sc.blockInfo.BlockNumber
			msBlock, _ := json.Marshal(sc.blockInfo)
			_, err = conn.Write([]byte(msBlock))
			if err != nil {
				logger.Error(err, "error accepted in the listener")
				break
			}
		}
		time.Sleep(time.Second)
	}
}
