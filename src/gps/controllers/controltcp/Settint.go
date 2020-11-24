package controltcp

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"goPanel/src/gps/coer/socket"
	"unsafe"
)

func SettingInit(cli unsafe.Pointer, message interface{}) {

}

func RegisterNode(cli unsafe.Pointer, message interface{}) {
	controlTcpCli := (*socket.Control)(unsafe.Pointer(cli))
	messBody, err := common.InterfaceByMapStr(message.(*socket.Message).Data)
	if err != nil {
		log.Error(err)
		return
	}

	controlTcpCli.Name = messBody["name"].(string)
	controlTcpCli.ClientId = messBody["uid"].(string)
}
