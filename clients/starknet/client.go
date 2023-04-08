package starknet_client

import (
	"context"

	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/ulerdogan/caigo-rpcv02/rpcv02"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

type starknetClient struct {
	Rpc *rpc.Provider
}

func NewStarknetClient(cnfg config.Config) Client {
	client, _ := ethrpc.DialContext(context.Background(), cnfg.RPCAddress)

	return &starknetClient{
		Rpc: rpc.NewProvider(client),
	}
}

func (c *starknetClient) Call(fc types.FunctionCall) ([]string, error) {
	return c.Rpc.Call(context.Background(), fc, rpcv02.WithBlockTag("latest"))
}

func (c *starknetClient) GetEvents(from, to uint64, address string, c_token *string, keys []string) ([]rpc.EmittedEvent, error) {
	output, err := c.Rpc.Events(context.Background(), rpc.EventsInput{
		FromBlock:         getBlockId(from),
		ToBlock:           getBlockId(to),
		Address:           getAddressHash(address),
		Keys:              keys,
		ContinuationToken: c_token,
		ChunkSize:         1024,
	})
	if err != nil {
		return nil, err
	}

	return output.Events, nil
}

func (c *starknetClient) LastBlock() (uint64, error) {
	return c.Rpc.BlockNumber(context.Background())
}
