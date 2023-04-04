package starknet_client

import (
	rpc "github.com/dontpanicdao/caigo/rpcv01"
	"github.com/dontpanicdao/caigo/types"
)

func getBlockId(number uint64) rpc.BlockID {
	return rpc.BlockID{
		Number: &number,
	}
}

func getAddressHash(address string) types.Hash {
	return types.HexToHash(address)
}