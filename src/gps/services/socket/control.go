package socket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

// 控制端
type ControlTcpManager struct {
	Client         map[*Control]bool // 连接控制端的客户端
	Broadcast      chan []byte       // 广播
	Register       chan *Control     // 注册
	UnRegister     chan *Control     // 卸载
	relayStartPort int               // 中继开始端口
	recoveryPort   []int             // 回收的端口
}

type Control struct {
	Conn  *net.Conn
	write chan []byte
}

var ControlManager = &ControlTcpManager{
	Client:         make(map[*Control]bool),
	Broadcast:      make(chan []byte, 1024),
	Register:       make(chan *Control),
	UnRegister:     make(chan *Control),
	relayStartPort: 10086,
}

func (cm *ControlTcpManager) Start() {
	controlListen, err := net.Listen("tcp", controlAddr)
	if err != nil {
		log.Error(err)
	}

	log.Info("控制端启动")

	go cm.Conn(controlListen)

	for {
		select {
		case cli := <-cm.Register:
			cm.Client[cli] = true
			go cli.Read()
			go cli.Send()
			msgJson, _ := json.Marshal(Message{
				Type:  0,
				Event: "connRelayByLocalSsh",
				Data:  "lallal",
				Code:  0,
			})

			cli.write <- msgJson
		case unRegister := <-cm.UnRegister:
			_ = (*unRegister.Conn).Close()
			delete(cm.Client, unRegister)
		case mess := <-cm.Broadcast:
			cm.SendAll(mess)
		}
	}
}

func (cm *ControlTcpManager) Conn(controlListen net.Listener) {
	for {
		controlConn, err := controlListen.Accept()
		if err != nil {
			log.Error(err)
		}

		log.Info("控制端有连接进来（新客户端）：", controlConn.RemoteAddr().String())

		client := Control{
			Conn:  &controlConn,
			write: make(chan []byte, 1024),
		}
		ControlManager.Register <- &client
	}
}

func (cm *ControlTcpManager) SendAll(message []byte) {
	if len(cm.Client) == 0 {
		return
	}

	for index, _ := range cm.Client {
		_, err := (*index.Conn).Write(message)

		if err != nil {
			log.Error(err)
		}
	}
}

func (c *Control) Read() {
	defer func() {
		if err := recover(); err != nil {
			ControlManager.UnRegister <- c
		}
	}()

	for {
		var data = make([]byte, 10240)
		size, err := (*c.Conn).Read(data)
		if err != nil || err == io.EOF {
			log.Info(err)
			ControlManager.UnRegister <- c
			break
		}
		data = data[:size]

		var dataMap map[string]interface{}
		err = json.Unmarshal(data, &dataMap)
		if err != nil {
			log.Info(err)
			continue
		}

		switch dataMap["event"] {
		case "init":
			serMess := map[string]string{
				"event": "createRelay",
				//"port":  strconv.Itoa(RelayPort()),
			}
			serMessJson, _ := json.Marshal(serMess)
			c.write <- serMessJson

			break
		}
	}
}

func (c *Control) Send() {
	defer func() {
		if err := recover(); err != nil {
			ControlManager.UnRegister <- c
		}
	}()

	for {
		select {
		case wr := <-c.write:
			_, err := (*c.Conn).Write(wr)
			if err != nil {
				log.Error("控制端发送消息失败！", err)
			}
		}

	}
}
