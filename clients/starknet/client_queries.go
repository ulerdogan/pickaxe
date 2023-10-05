package starknet_client

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

type Client interface {
	Call(fc rpc.FunctionCall) ([]*felt.Felt, error)
	GetEvents(from, to uint64, address string, c_token string, keys []string) ([]rpc.EmittedEvent, string, error)
	GetEventsWithID(from, to rpc.BlockID, address string, c_token string, keys []string) ([]rpc.EmittedEvent, string, error)
	LastBlock() (*rpc.BlockHashAndNumberOutput, error)
	NewDex(amm_id int) (Dex, error)
}
