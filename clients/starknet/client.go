package starknet_client

import (
	"context"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

func NewStarknetClient(cnfg config.Config) *starknetClient {
	client, _ := ethrpc.DialContext(context.Background(), cnfg.RPCAddress)

	return &starknetClient{
		Rpc: rpc.NewProvider(client),
	}
}
