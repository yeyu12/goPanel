package models

import (
	"github.com/go-xorm/xorm"
	"time"
)

// 要执行的命令
type CommandModel struct {
	Id           int64     `json:"id"`
	MachineId    int64     `json:"machine_id" xorm:"index notnull"` // 主机id
	Passwd       string    `json:"passwd"`                          // 要执行的主机密码，执行完成后，需要删除
	Flag         int       `json:"flag" xorm:"default(1)"`          // 执行方式，1立即执行，2定时(计划执行)执行
	Command      string    `json:"command"`                         // 要执行的命令
	ExecTime     time.Time `json:"exec_time"`                       // 执行时间
	PlanExecTime time.Time `json:"plan_exec_time"`                  // 计划执行时间
	IsExec       int       `json:"is_exec" xorm:"default(0)"`       // 是否执行（0未执行，1已执行
	ExecResult   string    `json:"exec_result"`                     // 执行结果
	CreateTime   time.Time `json:"create_time"`
	CreateUid    int64     `json:"create_uid"`
	UpdateTime   time.Time `json:"update_time"`
}

func (m *CommandModel) Add(db *xorm.Engine, data *CommandModel) (int64, error) {
	return db.InsertOne(data)
}
