package starknet_client

import (
	"github.com/dontpanicdao/caigo/types"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
)

type Client interface {
	Call(fc types.FunctionCall) ([]string, error)
	GetEvents(from, to uint64, address string, c_token *string, keys []string) ([]rpc.EmittedEvent, *string, error)
	LastBlock() (uint64, error)
	NewDex(amm_id int) (Dex, error)
}
