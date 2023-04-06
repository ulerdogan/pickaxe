package starknet_client

import (
	"github.com/dontpanicdao/caigo/types"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
)

func getBlockId(number uint64) rpc.BlockID {
	if number == 0 {
		return rpc.BlockID{
			Tag: "latest",
		}
	}

	return rpc.BlockID{
		Number: &number,
	}
}

func getAddressHash(address string) *types.Hash {
	if address == "" {
		return nil
	}

	h := types.HexToHash(address)
	return &h
}
