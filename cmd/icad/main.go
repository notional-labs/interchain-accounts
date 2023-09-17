package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/cosmos/interchain-accounts/v7/app"
	"github.com/cosmos/interchain-accounts/v7/cmd/icad/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "ICAD", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
