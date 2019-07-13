package main

import (
	"os"

	"kmodules.xyz/client-go/logs"
	"pharmer.dev/pre-k/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd(Version).Execute(); err != nil {
		os.Exit(1)
	}
}
