package starknet_client

import (
	"context"

	"github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

func (c *starknetClient) Call(fc types.FunctionCall) ([]string, error) {
	return c.Rpc.Call(context.Background(), fc, rpcv01.WithBlockTag("latest"))
}
