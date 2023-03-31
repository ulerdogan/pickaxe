package starknet_client

import (
	rpc "github.com/dontpanicdao/caigo/rpcv01"
)

type starknetClient struct {
	Rpc *rpc.Provider
}
