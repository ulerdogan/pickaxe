package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	//"github.com/ulerdogan/pickaxe/utils"
)

var testnet bool

func init() {
	//rootCmd.AddCommand(versionCmd)

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
		fmt.Println("app will be running soon in ", testnet)
	},
}

// var versionCmd = &cobra.Command{
// 	Use:   "version",
// 	Short: "Print the version number of Hugo",
// 	Long:  `All software has versions. This is Hugo's`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
// 	},
// }

func Execute() (err error) {
	err = rootCmd.Execute()
	return
}
