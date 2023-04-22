package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ulerdogan/pickaxe/socket"
)

var testnet bool

func init() {
	rootCmd.Flags().BoolVarP(&testnet, "testnet", "t", false, "use --testnet or -t to run on testnet")
}

var rootCmd = &cobra.Command{
	Use:   "psocket",
	Short: "starknet block finder",
	Long: `psocket listens new blocks and sends events to the tcp server
			when each block is created`,
	Version: "v0(dev)",
	Run: func(cmd *cobra.Command, args []string) {
		if testnet {
			socket.Init("app_test")
		} else {
			socket.Init("app")
		}
	},
}

func Execute() (err error) {
	err = rootCmd.Execute()
	return
}
