package main

import (
	"log"

	"github.com/ulerdogan/pickaxe/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("couldn't start the app: %v", err)
	}
}
