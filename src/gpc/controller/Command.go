package controller

import (
	"context"
	"encoding/json"
	"goPanel/src/gpc/service"
	"net"
)

func SendCommand(ctx context.Context, conn *net.TCPConn, message interface{}) {
	dataMap := message.(service.Message).Data.(map[string]interface{})
	dataJson, _ := json.Marshal(dataMap)
	var command service.Command
	_ = json.Unmarshal(dataJson, &command)

	service.ExecCommnadData = append(service.ExecCommnadData, command)
}
