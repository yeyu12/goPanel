package controltcp

import (
	"encoding/json"
	"goPanel/src/gps/coer/socket"
	"unsafe"
)

func Heartbeat(cli unsafe.Pointer, message interface{}) {
	controlTcpCli := (*socket.Control)(unsafe.Pointer(cli))
	req := socket.RequestWsMessage{
		Event: "heartbeat",
		Data:  nil,
	}
	reqJson, _ := json.Marshal(req)

	controlTcpCli.Write <- reqJson
}
