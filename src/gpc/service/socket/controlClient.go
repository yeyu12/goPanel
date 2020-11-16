package socket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"goPanel/src/gpc/router"
	"io"
)

func StartClientTcp(addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Info(err)
		}
	}()
	isReconnControlTcp = false
	clientConn, err := common.ConnTcp(addr)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		isReconnControlTcp = true
		clientConn.Close()
	}()

	for {
		var data = make([]byte, 10240)
		size, err := clientConn.Read(data)
		if err != nil || err == io.EOF {
			log.Error(err)
			break
		}
		data = data[:size]

		var message Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			log.Info(err)
			break
		}

		router.Route[message.Event](clientConn, message)
	}
}
