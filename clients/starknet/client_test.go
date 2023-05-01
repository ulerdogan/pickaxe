package starknet_client

import (
	"fmt"
	"testing"

	"github.com/dontpanicdao/caigo/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

func TestCall(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	paHash := types.HexToHash("0x00a144ef99419e4dbb3ef99bc2db894fbe7b4532ebed9592a407908727321fcf")
	r, err := c.Call(types.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: "getFee0",
		//Calldata:           []string{"1"},
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	dc := decimal.NewFromInt(types.HexToBN(r[0]).Int64()).Div(decimal.NewFromInt(10000)).String()
	fmt.Println(dc)

}

func TestGetEvents(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	events, _, err := c.GetEvents(
		29588,
		29589,
		"0x04d0390b777b424e43839cd1e744799f3de6c176c7e32c1812a41dbd9c19db6a",
		nil,
		[]string{"0xe14a408baf7f453312eec68e9b7d728ec5337fbdf671f917ee8c80f3255232"},
	)

	assert.Nil(t, err)
	fmt.Printf("Number of events found: %v\n", len(events))
}

func TestGetEventsWithID(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	from := uint64(29588)
	to := types.HexToHash("0x0072b6284d5003086dc23a568949f6e72129c3f594dbbee194ed150862e91dae")

	events, _, err := c.GetEventsWithID(
		rpc.BlockID{Number: &from},
		rpc.BlockID{Hash: &to},
		"0x04d0390b777b424e43839cd1e744799f3de6c176c7e32c1812a41dbd9c19db6a",
		nil,
		[]string{"0xe14a408baf7f453312eec68e9b7d728ec5337fbdf671f917ee8c80f3255232"},
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
