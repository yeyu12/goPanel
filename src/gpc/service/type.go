package service

import (
	"net"
	"time"
)

var ControlAddr string

type RequestWsMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Code  int32       `json:"code"`
}

type CommandService struct {
	Id           int64        `json:"id"`
	MachineId    string       `json:"machine_id"`     // 主机id
	Command      string       `json:"command"`        // 要执行的命令
	ExecTime     time.Time    `json:"exec_time"`      // 执行时间
	PlanExecTime time.Time    `json:"plan_exec_time"` // 计划执行时间
	IsExec       int          `json:"is_exec"`        // 是否执行（0未执行，1已执行
	ExecResult   string       `json:"exec_result"`    // 执行结果
	HandleTime   int64        `json:"handle_time"`    // 执行耗时
	IsLock       bool         // 锁
	Conn         *net.TCPConn // 连接信息
}
