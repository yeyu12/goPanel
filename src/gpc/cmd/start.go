package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	core_log "goPanel/src/core/log"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/core/signal"
	"goPanel/src/gpc/service"
	"goPanel/src/gpc/service/socket"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var StartCmd = cli.Command{
	Name:        "start",
	Usage:       "Start service gpc",
	Description: "Go panel client service startup",
	Action:      startRun,
	Flags:       []cli.Flag{},
}

func startRun(c *cli.Context) {
	go signal.HandleSignal()

	// TODO 重载服务的时候，需要重载配置

	_, _ = time.LoadLocation("Asia/Shanghai")

	conf := config.Conf.App
	core_log.Initialization(
		conf.LogOutputType,
		conf.Debug,
		conf.LogLevel,
	)
	core_log.LogSetOutput(conf.LogPath, conf.LogOutputFlag)

	// 写入pid文件
	if err := ioutil.WriteFile(config.GpcPidFileName, []byte(strconv.Itoa(os.Getpid())), 0755); err != nil {
		log.Panic(err)
	}

	service.ControlAddr = conf.ServerHost + ":" + conf.ServerPort
	socket.StartClientTcp()
}
