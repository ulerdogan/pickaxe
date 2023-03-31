package starknet_client

import (
	"context"

	rpc "github.com/dontpanicdao/caigo/rpcv01"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

func NewStarknetClient(cnfg config.Config) *starknetClient {
	client, _ := ethrpc.DialContext(context.Background(), cnfg.RPCAddress)

	return &starknetClient{
		Rpc: rpc.NewProvider(client),
	}
}
