package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func HandleSignal() {
	ch := make(chan os.Signal, 1)
	// 监听信号
	//signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	signal.Notify(ch, syscall.SIGUSR2)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGUSR2: // 重启

			return
		}
	}
}
