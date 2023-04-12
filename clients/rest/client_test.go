package rest_client

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c := NewRestClient()
	ticker := "ethereum"

	d := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", ticker)
	res, err := c.Get(d, nil)
	if err != nil {
		fmt.Println("coingecko query error")
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("coingecko query error")
		return
	}

	var price map[string]struct {
		Usd decimal.Decimal
	}

	json.Unmarshal(body, &price)
	mp := price[ticker].Usd

	assert.Equal(t, mp.GreaterThan(decimal.Zero), true)
	fmt.Println("ethereum price: ", mp)
}
