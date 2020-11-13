package socket

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func init() {
	controlListen, err := net.Listen("tcp", controlAddr)
	if err != nil {
		log.Error(err)
	}

	go func() {
		for {
			controlConn, err := controlListen.Accept()
			if err != nil {
				log.Error(err)
			}

			log.Info("控制启动，有连接进来：", controlConn.RemoteAddr().String())

			client := Control{
				Conn:  &controlConn,
				write: make(chan []byte, 1024),
			}
			ControlManager.Register <- &client
		}
	}()
}

func (cm *ControlTcpManager) Start() {
	for {
		select {
		case cli := <-cm.Register:
			cm.Client[cli] = true
			go cli.Read()
			go cli.Send()
		case unRegister := <-cm.UnRegister:
			_ = (*unRegister.Conn).Close()
			delete(cm.Client, unRegister)
		case mess := <-cm.Broadcast:
			cm.SendAll(mess)
		}
	}
}

func (cm *ControlTcpManager) SendAll(message []byte) {
	if len(cm.Client) == 0 {
		return
	}

	for index, _ := range cm.Client {
		_, err := (*index.Conn).Write(message)

		if err != nil {
			fmt.Println(err)
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
		data := make([]byte, 10240)
		size, err := (*c.Conn).Read(data)
		if err != nil {
			fmt.Println(err)
		}
		data = data[:size]

		var dataMap map[string]interface{}
		err = json.Unmarshal(data, &dataMap)
		if err != nil {
			fmt.Println("json解析失败", err)
		}

		switch dataMap["event"] {
		case "init":
			fmt.Println(dataMap["data"])
			// 那到数据后，解析出来，判断是否要创建隧道连接

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
				fmt.Println("控制端发送消息失败！", err)
			}
		}

	}
}
