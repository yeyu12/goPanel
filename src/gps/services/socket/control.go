package socket

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"goPanel/src/gps/router"
	"io"
	"net"
)

// 控制端
type ControlTcpManager struct {
	Clients        map[*Control]bool // 连接控制端的客户端
	Broadcast      chan []byte       // 广播
	Register       chan *Control     // 注册
	UnRegister     chan *Control     // 卸载
	relayStartPort int               // 中继开始端口
	recoveryPort   []int             // 回收的端口
}

type Control struct {
	Conn  *net.TCPConn
	write chan []byte
	Uuid  string
	Name  string
}

func (cm *ControlTcpManager) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", controlAddr)
	if err != nil {
		log.Panic(err)
		return
	}
	controlListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Panic(err)
		return
	}

	defer controlListen.Close()

	go cm.conn(controlListen)
	log.Info("控制端启动")

	for {
		select {
		case cli := <-cm.Register:
			cm.Clients[cli] = true
			go cli.read()
			go cli.send()
		case unRegister := <-cm.UnRegister:
			_ = (*unRegister.Conn).Close()
			delete(cm.Clients, unRegister)
		case mess := <-cm.Broadcast:
			cm.SendAll(mess)
		}
	}
}

func (cm *ControlTcpManager) conn(controlListen *net.TCPListener) {
	for {
		controlConn, err := controlListen.Accept()
		if err != nil {
			log.Error(err)
		}

		log.Info("控制端有连接进来（新客户端）：", controlConn.RemoteAddr().String())

		client := Control{
			Conn:  controlConn.(*net.TCPConn),
			write: make(chan []byte, 1024),
			Uuid:  uuid.NewV4().String(),
		}

		ControlManager.Register <- &client
	}
}

func (cm *ControlTcpManager) SendAll(message []byte) {
	if len(cm.Clients) == 0 {
		return
	}

	for index, _ := range cm.Clients {
		_, err := (*index.Conn).Write(message)

		if err != nil {
			log.Error(err)
			continue
		}
	}
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
			log.Error(err)
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

		if router.ControlTcpRoute[ret.Event] != nil {
			router.ControlTcpRoute[ret.Event](c.Conn, ret)
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
		case wr := <-c.write:
			_, err := (*c.Conn).Write(wr)
			if err != nil {
				log.Error("控制端发送消息失败！", err)
				return
			}
		}

	}
}
