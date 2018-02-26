package main

import (
	"fmt"
	"os"

	"github.com/paxautoma/operos/components/common"
	"github.com/paxautoma/operos/components/operosup/pkg/cmd"
)

func main() {
	common.SetupLogging()
	defer common.LogPanic()

	rootCmd := cmd.NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
