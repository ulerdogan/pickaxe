package starknet_client

import (
	"context"

	rpc "github.com/dontpanicdao/caigo/rpcv01"
)

func (c *starknetClient) Get(from uint64, address string, keys []string) ([]rpc.EmittedEvent, error) {
	output, err := c.Rpc.Events(context.Background(), rpc.EventFilter{
		FromBlock:  getBlockId(from),
		ToBlock:  getBlockId(from+1),
		Address:    getAddressHash(address),
		Keys:       keys,
		PageNumber: 0,
	})
	if err != nil {
		return nil, err
	}

	return output.Events, nil
}
