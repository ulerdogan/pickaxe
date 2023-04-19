package starknet_client

import (
	"context"
	"errors"

	"github.com/dontpanicdao/caigo/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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
	return c.Rpc.Call(context.Background(), fc, rpc.WithBlockTag("latest"))
}

func (c *starknetClient) GetEvents(from, to uint64, address string, c_token *string, keys []string) ([]rpc.EmittedEvent, *string, error) {
	output, err := c.Rpc.Events(context.Background(), rpc.EventsInput{
		FromBlock:         getBlockId(from),
		ToBlock:           getBlockId(to),
		Address:           getAddressHash(address),
		Keys:              keys,
		ContinuationToken: c_token,
		ChunkSize:         10,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.Events, output.ContinuationToken, nil
}

func (c *starknetClient) LastBlock() (uint64, error) {
	return c.Rpc.BlockNumber(context.Background())
}

func (c *starknetClient) NewDex(amm_id int) (Dex, error) {
	switch amm_id {
	case 1:
		return newMyswap(), nil
	case 2:
		return newJediswap(), nil
	case 3:
		return newSwap10k(), nil
	}

	return nil, errors.New("cannot find the dex")
}
