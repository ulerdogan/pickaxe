package starknet_client

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
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

func GetAddressFelt(address string) *felt.Felt {
	if address == "" {
		return nil
	}

	h, _ := utils.HexToFelt(address)
	return h
}

func getStrBigIntFelt(number string) *felt.Felt {
	f, err := new(felt.Felt).SetString(number)
	if err != nil {
		return nil
	}

	return f
}

// Get a unique pool identifier for the Ekubo pools with given data:
// TokenA, TokenB, Fee, TickSpacing
func GetUniqueEkuboHash(i, j, k, q string) string {
	combined := i + j + k + q
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func GetAdressFormatFromFelt(fl *felt.Felt) string {
	b, _ := types.HexToBytes(fl.String())
	return "0x" + strings.Repeat("0", 64-len(hex.EncodeToString(b))) + hex.EncodeToString(b)
}

func GetAdressFormatFromStr(s string) string {
	hx, _ := hex.DecodeString(s[2:])
	return "0x" + strings.Repeat("0", 64-len(hex.EncodeToString(hx))) + hex.EncodeToString(hx)
}
