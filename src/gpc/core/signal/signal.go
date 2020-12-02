package signal

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/gpc/service/socket"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGUSR2: // 重启
			log.Info("reboot service!")
			socket.Cancel()
			log.Info("reboot service success!")
		}
	}
}
