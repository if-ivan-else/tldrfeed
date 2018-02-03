package main

import (
	"log"

	"github.com/if-ivan-else/tldrfeed/cmd/tldrfeed/app"
)

// main entry point of the program
func main() {
	err := app.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
