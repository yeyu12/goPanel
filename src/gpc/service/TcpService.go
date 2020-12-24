package service

import (
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
		if _, err = s.Conn.Write([]byte(item)); err != nil {
			return err
		}

		time.Sleep(time.Microsecond * 100)
	}

	return nil
}

// 接收
