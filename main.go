package main

import (
	"fmt"
	"os"
	"github.com/mcnull/qai/app"
	"github.com/mcnull/qai/shared/utils"
)

func main() {

	utils.LoadEnvFile()

	app := app.NewApp()
	c, err := app.Init(os.Args)

	if err != nil {
		fmt.Printf("Error initializing: %v\n", err)
		os.Exit(1)
	}

	if !c {
		os.Exit(0)
	}

	// Run the app
	err = app.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
