package starknet_client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

func TestGet(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	events, err := c.Get(
		29588,
		29589,
		"0x04d0390b777b424e43839cd1e744799f3de6c176c7e32c1812a41dbd9c19db6a",
		nil,
		[]string{"0xe14a408baf7f453312eec68e9b7d728ec5337fbdf671f917ee8c80f3255232"},
	)

	assert.Nil(t, err)
	fmt.Printf("Number of events found: %v\n", len(events))
}
