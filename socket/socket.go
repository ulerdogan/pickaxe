package socket

import (
	"net"
	"sync"
	"time"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/go-co-op/gocron"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

type socket struct {
	listener net.Listener
	client   starknet.Client

	config config.Config

	scheduler *gocron.Scheduler
	scMutex   *sync.Mutex
	blockInfo *rpc.BlockHashAndNumberOutput
}

func NewSocket(ls net.Listener, cl starknet.Client, cnfg config.Config) *socket {
	return &socket{
		listener:  ls,
		client:    cl,
		config:    cnfg,
		scheduler: gocron.NewScheduler(time.UTC),
		scMutex:   &sync.Mutex{},
		blockInfo: nil,
	}
}
