package starknet_client

import (
	"context"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
)

func (c *starknetClient) Get(from, to uint64, address string, c_token *string, keys []string) ([]rpc.EmittedEvent, error) {
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
