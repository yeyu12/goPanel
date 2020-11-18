package router

import (
	"errors"
	"unsafe"
)

type handle func(cli unsafe.Pointer, message interface{})

var controlTcpRoute = make(map[string]handle)

func AddRoute(path string, handleFunc handle) {
	controlTcpRoute[path] = handleFunc
}

func HandleRoute(path string, cli unsafe.Pointer, message interface{}) error {
	if controlTcpRoute[path] != nil {
		controlTcpRoute[path](cli, message)
		return nil
	}

	return errors.New("controlTcp路由未找到！")
}
