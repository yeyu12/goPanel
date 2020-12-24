package socket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"goPanel/src/core/tcp_package"
	"goPanel/src/gps/coer/router"
	"io"
	"net"
	"time"
	"unsafe"
)

type Control struct {
	Conn       *net.TCPConn
	Write      chan []byte
	Uuid       string
	Name       string
	ClientId   string                                          // 客户端uid
	SystemType string                                          // 系统信息
	TcpBody    map[int64]map[int64]*tcp_package.PackageContent // 消息分包包体
}

func (c *Control) read() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			ControlManager.UnRegister <- c
		}
	}()

	tcpPackageObj := tcp_package.NewTcpPackage(constants.DEFAULT_SUBPACKAGE, time.Now().UnixNano())

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
		// 拆包
		unPackingData, err := tcpPackageObj.TcpUnPacking(data)
		if err != nil {
			log.Info(err)
			continue
		}

		if _, ok := c.TcpBody[unPackingData.PackageId]; !ok {
			c.TcpBody[unPackingData.PackageId] = make(map[int64]*tcp_package.PackageContent)
		}
		c.TcpBody[unPackingData.PackageId][unPackingData.PackageIndex] = unPackingData

		body, err := tcpPackageObj.TcpJoinPackage(c.TcpBody[unPackingData.PackageId])
		if err != nil {
			log.Debug(err)
			continue
		}

		var ret Message
		err = json.Unmarshal(body, &ret)
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
