package starknet_client

import (
	"crypto/sha256"
	"encoding/hex"

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

func GetUniqueEkuboHash(i, j, k, q string) string {
	combined := i + j + k + q
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}
