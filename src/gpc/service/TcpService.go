package service

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"goPanel/src/core/tcp_package"
	"net"
	"time"
)

func NewTcpService(conn *net.TCPConn) *TcpService {
	return &TcpService{
		Conn: conn,
	}
}

// 发送
func (s *TcpService) Send(data []byte) error {
	subpackageData, err := tcp_package.NewTcpPackage(constants.DEFAULT_SUBPACKAGE, time.Now().UnixNano()).TcpSubpackage(string(data))
	if err != nil {
		return err
	}

	for _, item := range subpackageData {
		size, err := s.Conn.Write([]byte(item))
		if err != nil {
			return err
		}

		log.Error(size, err)

		//time.Sleep(time.Microsecond * 100)
	}

	return nil
}
