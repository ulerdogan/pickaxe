package cmd

import (
	"github.com/spf13/cobra"
	init_db "github.com/ulerdogan/pickaxe/db/init"
	
	"github.com/ulerdogan/pickaxe/indexer"
)

var testnet bool

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.Flags().BoolVarP(&testnet, "testnet", "t", false, "use --testnet or -t to run on testnet")
}

var rootCmd = &cobra.Command{
	Use:   "pickaxe",
	Short: "sister of shovel",
	Long: `pickaxe is AMM pool indexer for starknet defi ecosystem.
				  it indexes choosen starknet dexes amm pools to prepare
				  data flow for https://fibrous.finance.`,
	Version: "v0(dev)",
	Run: func(cmd *cobra.Command, args []string) {
		if testnet {
			indexer.Init("app_test")
		} else {
			indexer.Init("app")
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init token-pool-amm items for the db",
	Long: `pickaxe is needed to be initialized with the first state
				  of the amm pools, tokens, and amms
				  which is in the db\init folder and the architecture is
				  designed for the allowing first runs efficiently.`,
	Run: func(cmd *cobra.Command, args []string) {
		init_db.Init()
	},
}

func Execute() (err error) {
	err = rootCmd.Execute()
	return
}
