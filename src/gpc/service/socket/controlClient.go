package socket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/router"
	"io"
	"net"
	"time"
)

func StartClientTcp() {
	for {
		if isReconnControlTcp {
			handelConnControlTcp()
			log.Info("重连重试中！")
		}

		time.Sleep(time.Second * time.Duration(config.Conf.App.ControlReconnTcpTime))
	}
}

// 心跳
func heartbeat(conn *net.TCPConn) {
	for {
		time.Sleep(time.Second * time.Duration(config.Conf.App.ControlHeartbeatTime))

		write := RequestWsMessage{
			Event: "heartbeat",
			Data:  nil,
		}
		writeJson, err := json.Marshal(write)
		if err != nil {
			continue
		}

		log.Info("正在执行控制端心跳包")
		if _, err = conn.Write(writeJson); err != nil {
			log.Error(err)
			return
		}
	}
}

// 注册本机数据
func registerLocalData(conn *net.TCPConn) {
	localComputerData := map[string]string{
		"name": config.Conf.App.LocalName,
	}
	write := RequestWsMessage{
		Event: "local_register",
		Data:  localComputerData,
	}
	writeJson, err := json.Marshal(write)
	if err != nil {
		log.Error(err)
		return
	}
	if _, err = conn.Write(writeJson); err != nil {
		log.Error(err)
		return
	}
}

func handelConnControlTcp() {
	defer func() {
		isReconnControlTcp = true
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	isReconnControlTcp = false
	conn, err := common.ConnTcp(ControlAddr)
	if err != nil {
		log.Error(err)
		return
	}

	defer func() {
		conn.Close()
	}()

	registerLocalData(conn)
	go heartbeat(conn)

	for {
		var data = make([]byte, 10240)
		size, err := conn.Read(data)
		if err != nil || err == io.EOF {
			log.Error(err)
			break
		}
		data = data[:size]

		var message Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			log.Info(err)
			continue
		}

		router.Route[message.Event](conn, message)
	}
}
