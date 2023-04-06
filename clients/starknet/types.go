package starknet_client

import (
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
)

type starknetClient struct {
	Rpc *rpc.Provider
}
