package socket

import (
	"net"

	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func Init(environment string) {
	// load app configs
	cnfg, err := config.LoadConfig(environment, ".")
	if err != nil {
		logger.Error(err, "cannot load config for psocket: "+environment)
		return
	}
	logger.Info("config loaded for psocket: " + environment)

	client := starknet.NewStarknetClient(cnfg)

	ls, err := net.Listen("tcp", cnfg.SocketAddress)
	if err != nil {
		logger.Error(err, "cannot listen to: "+cnfg.SocketAddress)
		return
	}
	
	logger.Info("socket server will be running at: " + cnfg.SocketAddress)

	sc := NewSocket(ls, client, cnfg)

	// setup and run jobs
	setupJobs(sc)

	go sc.Sync()
	sc.scheduler.StartBlocking()
}
