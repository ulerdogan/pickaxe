package socket

import (
	"strconv"
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
			str := strconv.Itoa(int(lastSent))
			_, err = conn.Write([]byte(str))
			if err != nil {
				logger.Error(err, "error accepted in the listener")
				break
			}
		}
		time.Sleep(time.Second)
	}
}
