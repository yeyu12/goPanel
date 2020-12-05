package socket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gps/coer/router"
	"io"
	"net"
	"unsafe"
)

type Control struct {
	Conn       *net.TCPConn
	Write      chan []byte
	Uuid       string
	Name       string
	ClientId   string // 客户端uid
	SystemType string // 系统信息
}

func (c *Control) read() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			ControlManager.UnRegister <- c
		}
	}()

	for {
		var data = make([]byte, 10240)
		size, err := (*c.Conn).Read(data)
		if err != nil || err == io.EOF {
			if err != io.EOF {
				log.Error(err)
			}

			ControlManager.UnRegister <- c
			break
		}
		data = data[:size]

		var ret Message
		err = json.Unmarshal(data, &ret)
		if err != nil {
			log.Error(err)
			continue
		}

		if err = router.HandleRoute(ret.Event, unsafe.Pointer(c), &ret); err != nil {
			log.Error(err)
		}
	}
}

func (c *Control) send() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			ControlManager.UnRegister <- c
		}
	}()

	for {
		select {
		case wr := <-c.Write:
			_, err := (*c.Conn).Write(wr)
			if err != nil {
				log.Error("控制端发送消息失败！", err)
				return
			}
		}

	}
}
