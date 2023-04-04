package starknet_client

import (
	"fmt"
	"testing"

	"github.com/dontpanicdao/caigo/types"
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
