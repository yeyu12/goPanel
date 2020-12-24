package socket

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"goPanel/src/core/tcp_package"
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
	log.Info("控制端TCP启动：" + controlListen.Addr().String())

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
			Conn:    controlConn.(*net.TCPConn),
			Write:   make(chan []byte, 1024),
			Uuid:    uuid.NewV4().String(),
			TcpBody: make(map[int64]map[int64]*tcp_package.PackageContent),
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

// 客户端id获取客户端连接信息
func (cm *ControlTcpManager) FindClientIdByClientConn(clientId string) *Control {
	for index, _ := range ControlManager.Clients {
		if index.ClientId == clientId {
			return index
		}
	}

	return nil
}

func (cm *ControlTcpManager) PushRecoveryPort(port int) {
	ControlManager.recoveryPort = append(ControlManager.recoveryPort, port)
}
