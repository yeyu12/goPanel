package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"goPanel/src/common"
	"goPanel/src/constants"
	"io/ioutil"
	"os/exec"
)

var KillCmd = cli.Command{
	Name:        "kill",
	Usage:       "Kill service gpc",
	Description: "Go panel client service startup",
	Action:      killRun,
	Flags:       []cli.Flag{},
}

func killRun(c *cli.Context) {
	pidStr, err := ioutil.ReadFile(common.GetCurrentDir() + constants.GPC_PID_PATH + constants.PID_FILENAME)
	if err != nil {
		log.Error(err)
	}

	cmd := exec.Command("kill", "-9", string(pidStr))
	_ = cmd.Run()
}
