package init_db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
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
	Address     string `json:"address"`
	TokenA      string `json:"token_a"`
	TokenB      string `json:"token_b"`
	ExtraData   string `json:"extra_data,omitempty"`
	AmmID       int    `json:"amm_id"`
	Fee         string `json:"fee,omitempty"`
	TickSpacing string `json:"tick_spacing,omitempty"`
}

var (
	tokens    []token
	amms      []amm
	pools     []pool
	ekuboFees map[string]string = make(map[string]string)
)

func initStates(initPath string) error {

	//tokensFile, err := os.Open("./init/states/tokens.json")
	tokensFile, err := os.Open(fmt.Sprintf("%s/%s.json", initPath, "tokens"))
	if err != nil {
		logger.Error(err, "cannot get tokens json file")
		return err
	}
	defer tokensFile.Close()

	ammsFile, err := os.Open(fmt.Sprintf("%s/%s.json", initPath, "amms"))
	if err != nil {
		logger.Error(err, "cannot get amms json file")
		return err
	}
	defer ammsFile.Close()

	poolsFile, err := os.Open(fmt.Sprintf("%s/%s.json", initPath, "pools"))
	if err != nil {
		logger.Error(err, "cannot get pools json file")
		return err
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

	return nil
}

func initTokensToDB(store db.Store, c starknet.Client) {
	succesfulls := 0

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
			continue
		}

		succesfulls++
	}
	logger.Info("in " + strconv.Itoa(len(tokens)) + " tokens, " + strconv.Itoa(succesfulls) + " is created")
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

	for _, p := range pools {
		ta, _ := store.GetTokenBySymbol(context.Background(), p.TokenA)
		tb, _ := store.GetTokenBySymbol(context.Background(), p.TokenB)

		// sort pools
		if ta.Address > tb.Address {
			ta, tb = tb, ta
		}

		pool, err := store.CreatePool(context.Background(), db.CreatePoolParams{
			Address: p.Address,
			TokenA:  ta.Address,
			TokenB:  tb.Address,
			AmmID:   int64(p.AmmID),
		})
		if err != nil {
			logger.Error(err, "cannot create pool: "+p.Address)
			continue
		}

		if p.ExtraData != "" && p.AmmID != 5 {
			store.UpdatePoolExtraData(context.Background(), db.UpdatePoolExtraDataParams{
				PoolID:    pool.PoolID,
				ExtraData: sql.NullString{String: p.ExtraData, Valid: true},
			})
		}

		if p.Fee != "" && p.TickSpacing != "" {
			ekuboData := starknet.EkuboData{
				TickSpacing:  p.TickSpacing,
				KeyExtension: "0",
			}

			jsonBytes, _ := json.Marshal(ekuboData)
			pool, err = store.UpdatePoolGeneralExtraData(context.Background(), db.UpdatePoolGeneralExtraDataParams{
				PoolID:           pool.PoolID,
				GeneralExtraData: sql.NullString{String: string(jsonBytes), Valid: true},
			})
			if err != nil {
				logger.Error(err, "cannot create pool: "+p.Address)
				continue
			}

			ekuboExtraData := starknet.GetUniqueEkuboHash(p.TokenA, p.TokenB, p.Fee, p.TickSpacing)
			pool, err = store.UpdatePoolExtraData(context.Background(), db.UpdatePoolExtraDataParams{
				PoolID:    pool.PoolID,
				ExtraData: sql.NullString{String: ekuboExtraData, Valid: true},
			})

			if err != nil {
				logger.Error(err, "cannot create pool: "+p.Address)
				continue
			}

			ekuboFees[ekuboExtraData] = p.Fee
		}

		succesfulls++
	}
	logger.Info("in " + strconv.Itoa(len(pools)) + " pools, " + strconv.Itoa(succesfulls) + " is created")
}

func getTokenDecimal(c starknet.Client, address string) (*int, error) {
	paHash := starknet.GetAddressFelt(address)
	r, err := c.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("decimals"),
		Calldata:           []*felt.Felt{},
	})

	if err != nil {
		logger.Error(err, "cannot get token decimal: "+address)
		return nil, err
	}

	decimal := int(r[0].BigInt(new(big.Int)).Int64())
	return &decimal, nil
}

func Init(cnfg config.Config, store db.Store, client starknet.Client) {
	if err := initStates(cnfg.InitPath); err != nil {
		logger.Error(err, "cannot init states")
		return
	}

	logger.Info("first state initialization runned")

	initAmmsToDB(store)
	initTokensToDB(store, client)
	initPoolsToDB(store)

	pools, _ := store.GetAllPools(context.Background())

	const numWorkers = 10
	jobs := make(chan db.Pool, len(pools))
	results := make(chan bool, len(pools))

	block, err := client.LastBlock()
	if err != nil {
		logger.Error(err, "cannot get the last block")
	}

	for w := 0; w < numWorkers; w++ {
		go func(jobs chan db.Pool, results chan bool) {
			syncPoolFromFnConc(jobs, results, block.BlockNumber, store, client)
		}(jobs, results)
	}

	for _, pool := range pools {
		jobs <- pool
	}
	close(jobs)

	var s int
	for res := 0; res < len(pools); res++ {
		if <-results {
			s++
		}
	}

	logger.Info("in " + strconv.Itoa(len(pools)) + " pools, " + strconv.Itoa(s) + " is synced")
}

func syncPoolFromFnConc(jobs <-chan db.Pool, results chan<- bool, lastBlock uint64, store db.Store, client starknet.Client) {
	for pool := range jobs {
		dex, _ := client.NewDex(int(pool.AmmID))

		syncFeeInfo := starknet.PoolInfo{
			Address:   pool.Address,
			ExtraData: pool.ExtraData.String,
			Fee:       ekuboFees[pool.ExtraData.String],
		}

		err := dex.SyncFee(syncFeeInfo, store, client)
		if err != nil {
			logger.Error(err, "sync fee error: "+pool.Address)
			results <- false
			continue
		}

		err = dex.SyncPoolFromFn(starknet.PoolInfo{
			Address:   pool.Address,
			ExtraData: pool.ExtraData.String,
			Block:     big.NewInt(int64(lastBlock)),
		}, store, client)
		if err != nil {
			logger.Error(err, "sync pool error: "+pool.Address)
			results <- false
			continue
		}

		results <- true
	}
}
