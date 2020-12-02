package main

import (
	"github.com/urfave/cli"
	"goPanel/src/constants"
	"goPanel/src/gpc/cmd"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gpc"
	app.Version = constants.GPC_VERSION
	app.Usage = "Control panel client based on go"
	app.Commands = []cli.Command{
		cmd.StartCmd,
		cmd.KillCmd,
	}

	if err := app.Run(os.Args); err != nil {
		log.Panic(err)
	}
}
