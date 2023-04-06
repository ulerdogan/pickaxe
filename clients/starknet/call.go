package starknet_client

import (
	"context"

	"github.com/dontpanicdao/caigo/types"
	"github.com/ulerdogan/caigo-rpcv02/rpcv02"
)

func (c *starknetClient) Call(fc types.FunctionCall) ([]string, error) {
	return c.Rpc.Call(context.Background(), fc, rpcv02.WithBlockTag("latest"))
}
