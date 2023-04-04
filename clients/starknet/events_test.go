package starknet_client

import (
	"fmt"
	"testing"

	config "github.com/ulerdogan/pickaxe/utils/config"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cnfg, _ := config.LoadConfig("app", "../..")
	c := NewStarknetClient(cnfg)

	events, err := c.Get(
		29588,
		"0x04d0390b777b424e43839cd1e744799f3de6c176c7e32c1812a41dbd9c19db6a",
		[]string{"0xe14a408baf7f453312eec68e9b7d728ec5337fbdf671f917ee8c80f3255232"},
	)

	assert.Nil(t, err)
	fmt.Println(events)
}
