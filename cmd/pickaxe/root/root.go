package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ulerdogan/pickaxe/indexer"
)

var testnet bool

func init() {
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

func Execute() (err error) {
	err = rootCmd.Execute()
	return
}
