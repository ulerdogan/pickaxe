package starknet_client

import (
	"context"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	utils "github.com/NethermindEth/starknet.go/utils"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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

func (c *starknetClient) Call(fc rpc.FunctionCall) ([]*felt.Felt, error) {
	return c.Rpc.Call(context.Background(), fc, rpc.WithBlockTag("pending"))
}

func (c *starknetClient) GetEvents(from, to uint64, address string, c_token string, keys []string) ([]rpc.EmittedEvent, string, error) {
	var felt_keys [][]*felt.Felt
	felt_keys = append(felt_keys, []*felt.Felt{})

	for i := range keys {
		feltKey, _ := utils.HexToFelt(keys[i])
		felt_keys[0] = append(felt_keys[0], feltKey)
	}

	output, err := c.Rpc.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: getBlockId(from),
			ToBlock:   getBlockId(to),
			Address:   GetAddressFelt(address),
			Keys:      felt_keys,
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ContinuationToken: "",
			ChunkSize:         1024,
		},
	})
	if err != nil {
		return nil, "", err
	}

	return output.Events, output.ContinuationToken, nil
}

func (c *starknetClient) GetEventsWithID(from, to rpc.BlockID, address string, c_token string, keys []string) ([]rpc.EmittedEvent, string, error) {
	var felt_keys [][]*felt.Felt
	felt_keys = append(felt_keys, []*felt.Felt{})

	for i := range keys {
		feltKey, _ := utils.HexToFelt(keys[i])
		felt_keys[0] = append(felt_keys[0], feltKey)
	}

	output, err := c.Rpc.Events(context.Background(), rpc.EventsInput{
		EventFilter: rpc.EventFilter{
			FromBlock: from,
			ToBlock:   to,
			Address:   GetAddressFelt(address),
			Keys:      felt_keys,
		},
		ResultPageRequest: rpc.ResultPageRequest{
			ContinuationToken: c_token,
			ChunkSize:         1024,
		},
	})
	if err != nil {
		return nil, "", err
	}

	return output.Events, output.ContinuationToken, nil
}

func (c *starknetClient) LastBlock() (*rpc.BlockHashAndNumberOutput, error) {
	return c.Rpc.BlockHashAndNumber(context.Background())
}

func (c *starknetClient) NewDex(amm_id int) (Dex, error) {
	switch amm_id {
	case 1:
		return newMyswap(), nil
	case 2:
		return newJediswap(), nil
	case 3:
		return newSwap10k(), nil
	case 4:
		return newSithswap(), nil
	case 5:
		return newEkubo(), nil
	}

	return nil, errors.New("cannot find the dex")
}
