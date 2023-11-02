package starknet_client

import (
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

func TestCall(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	address := "0x00a144ef99419e4dbb3ef99bc2db894fbe7b4532ebed9592a407908727321fcf"
	r, err := c.Call(rpc.FunctionCall{
		ContractAddress:    GetAddressFelt(address),
		EntryPointSelector: types.GetSelectorFromNameFelt("getFee0"),
		//Calldata:           []string{"1"},
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	feltInBigInt, _ := utils.FeltToBigInt(r[0])
	dc := decimal.NewFromInt(feltInBigInt.Int64()).Div(decimal.NewFromInt(10000)).String()
	fmt.Println(dc)
}

func TestGetEvents(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	events, _, err := c.GetEvents(
		29588,
		29589,
		"0x04d0390b777b424e43839cd1e744799f3de6c176c7e32c1812a41dbd9c19db6a",
		"",
		[]string{"0xe14a408baf7f453312eec68e9b7d728ec5337fbdf671f917ee8c80f3255232"},
	)

	assert.Nil(t, err)
	fmt.Printf("Number of events found: %v\n", len(events))
}

func TestGetEventsWithID(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	from := uint64(325998)
	to := GetAddressFelt("0x038e1a8e97ca5729650d44858e2c8597b4f1ad5d046f73cc3d2b394b2e3fc2a2")

	events, _, err := c.GetEventsWithID(
		rpc.BlockID{Number: &from},
		rpc.BlockID{Hash: to},
		"",
		"",
		[]string{"0xe14a408baf7f453312eec68e9b7d728ec5337fbdf671f917ee8c80f3255232", "0x157717768aca88da4ac4279765f09f4d0151823d573537fbbeb950cdbd9a870"},
	)

	assert.Nil(t, err)
	fmt.Printf("Number of events found: %v\n", len(events))
}

func TestLastBlock(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	block, err := c.LastBlock()
	assert.Nil(t, err)
	fmt.Printf("Block number: %v\n", block.BlockNumber)
	fmt.Printf("Block hash: %v\n", block.BlockHash)
}
