package main

import (
	"github.com/ulerdogan/pickaxe/cmd"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logger.Error(err, "couldn't start the app")
	}
}
