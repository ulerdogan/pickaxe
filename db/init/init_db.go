package init_db

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/dontpanicdao/caigo/types"
	"github.com/shopspring/decimal"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

type token struct {
	Symbol  string `json:"symbol"`
	Name    string `json:"name"`
	Ticker  string `json:"ticker"`
	Address string `json:"address"`
	Base    bool   `json:"base"`
	Native  bool   `json:"native"`
}

type amm struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	RouterAddress string `json:"router_address"`
	Key           string `json:"key"`
	AlgorithmType string `json:"algorithm_type"`
}

type pool struct {
	Address   string `json:"address"`
	TokenA    string `json:"token_a"`
	TokenB    string `json:"token_b"`
	ExtraData string `json:"extra_data,omitempty"`
	AmmID     int    `json:"amm_id"`
}

var (
	tokens []token
	amms   []amm
	pools  []pool
)

func init() {
	// FIXME: fix for to be runned from the init page
	tokensFile, err := os.Open("./db/init/states/tokens.json")
	if err != nil {
		logger.Error(err, "cannot get tokens json file")
		return
	}
	defer tokensFile.Close()

	ammsFile, err := os.Open("./db/init/states/amm.json")
	if err != nil {
		logger.Error(err, "cannot get amms json file")
		return
	}
	defer ammsFile.Close()

	poolsFile, err := os.Open("./db/init/states/pools.json")
	if err != nil {
		logger.Error(err, "cannot get pools json file")
		return
	}
	defer ammsFile.Close()

	var resultTokens map[string][]token
	var resultAmms map[string][]amm
	var resultPools map[string][]pool

	byteValue, _ := io.ReadAll(tokensFile)
	json.Unmarshal(byteValue, &resultTokens)
	byteValue, _ = io.ReadAll(ammsFile)
	json.Unmarshal(byteValue, &resultAmms)
	byteValue, _ = io.ReadAll(poolsFile)
	json.Unmarshal(byteValue, &resultPools)

	tokens = resultTokens["tokens"]
	amms = resultAmms["amms"]
	pools = resultPools["pools"]
}

func initTokensToDB(store db.Store, c starknet.Client) {
	succesfulls := 0
	bErr := 0

	for _, t := range tokens {
		decimal, err := getTokenDecimal(c, t.Address)
		if err != nil {
			d := 18
			decimal = &d
		}
		_, err = store.CreateToken(context.Background(), db.CreateTokenParams{
			Address:  t.Address,
			Name:     t.Name,
			Symbol:   t.Symbol,
			Decimals: int32(*decimal),
			Ticker:   t.Ticker,
		})
		if err != nil {
			logger.Error(err, "cannot create token: "+t.Symbol)
			continue
		}

		_, err = store.UpdateBaseNativeStatus(context.Background(), db.UpdateBaseNativeStatusParams{
			Address: t.Address,
			Base:    t.Base,
			Native:  t.Native,
		})
		if err != nil {
			logger.Error(err, "cannot update base native status for: "+t.Symbol)
			bErr++
			continue
		}

		succesfulls++
	}
	logger.Info("in " + strconv.Itoa(len(tokens)) + " tokens, " + strconv.Itoa(succesfulls) + " is created and " + strconv.Itoa(bErr) + " is not updated with base and native status")
}

func initAmmsToDB(store db.Store) {
	succesfulls := 0

	for _, a := range amms {
		_, err := store.CreateAmm(context.Background(), db.CreateAmmParams{
			DexName:       a.Name,
			RouterAddress: a.RouterAddress,
			Key:           a.Key,
			AlgorithmType: a.AlgorithmType,
		})
		if err != nil {
			logger.Error(err, "cannot create amm: "+a.Name)
			continue
		}
		succesfulls++
	}
	logger.Info("in " + strconv.Itoa(len(amms)) + " amms, " + strconv.Itoa(succesfulls) + " is created")
}

func initPoolsToDB(store db.Store) {
	succesfulls := 0
	dc03, _ := decimal.NewFromString("0.3")

	for _, p := range pools {
		ta, _ := store.GetTokenBySymbol(context.Background(), p.TokenA)
		tb, _ := store.GetTokenBySymbol(context.Background(), p.TokenB)

		pool, err := store.CreatePool(context.Background(), db.CreatePoolParams{
			Address: p.Address,
			TokenA:  ta.Address,
			TokenB:  tb.Address,
			AmmID:   int64(p.AmmID),
			Fee:     dc03.String(),
		})
		if err != nil {
			logger.Error(err, "cannot create pool: "+p.Address)
			continue
		}

		if p.ExtraData != "" {
			store.UpdatePoolExtraData(context.Background(), db.UpdatePoolExtraDataParams{
				PoolID:    pool.PoolID,
				ExtraData: sql.NullString{String: p.ExtraData, Valid: true},
			})
		}

		succesfulls++
	}
	logger.Info("in " + strconv.Itoa(len(pools)) + " pools, " + strconv.Itoa(succesfulls) + " is created")
}

func getTokenDecimal(c starknet.Client, address string) (*int, error) {
	paHash := types.HexToHash(address)
	r, err := c.Call(types.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: "decimals",
		Calldata:           []string{},
	})

	if err != nil {
		logger.Error(err, "cannot get token decimal: "+address)
		return nil, err
	}

	decimal := int(types.HexToBN(r[0]).Int64())
	return &decimal, nil
}

func Init(cnfg config.Config, store db.Store, client starknet.Client) {
	//TODO: add initial reserves logic
	logger.Info("first state initialization runned")

	initAmmsToDB(store)
	initTokensToDB(store, client)
	initPoolsToDB(store)
}
