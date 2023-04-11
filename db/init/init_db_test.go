package init_db

import (	
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert.Equal(t, len(tokens), 5)
	assert.Equal(t, len(amms), 3)
	assert.Greater(t, len(pools), 0)
}

// func TestGetTokenDecimals(t *testing.T) {
// 	cnfg, _ := config.LoadConfig("app", "../..")
// 	client := starknet.NewStarknetClient(cnfg)

// 	i, err := getTokenDecimal(client, tokens[0].Address)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(*i)
// }
