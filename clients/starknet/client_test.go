package starknet_client

import (
	"fmt"
	"testing"

	"github.com/dontpanicdao/caigo/types"
	"github.com/stretchr/testify/assert"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

func TestCall(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	paHash := types.HexToHash("0x010884171baf1914edc28d7afb619b40a4051cfae78a094a55d230f19e944a28")
	r, err := c.Call(types.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: "get_pool",
		Calldata:           []string{"1"},
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
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

func TestLastBlock(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	block, err := c.LastBlock()
	assert.Nil(t, err)
	fmt.Printf("Block number: %v\n", block.BlockNumber)
	fmt.Printf("Block hash: %v\n", block.BlockHash)
}
