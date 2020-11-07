package validations

type MachineAdd struct {
	Flag int `json:"flag" validate:"required" label:"类型"`
}

type MachineDel struct {
	Flag int   `json:"flag" validate:"required" label:"类型"`
	Id   int64 `json:"id" validate:"required" label:"ID"`
}

type MachineAddDir struct {
	Id   int    `json:"id" label:"ID"`
	Name string `json:"name" validate:"required" label:"目录名称"`
}

type MachineAddComputer struct {
	Id             int64  `json:"id" label:"ID"`
	MachineGroupId int64  `json:"machine_group_id" validate:"min=0" label:"目录ID"`
	Name           string `json:"name" label:"名称"`
	Host           string `json:"host" validate:"required" label:"地址"`
	User           string `json:"user" label:"用户名"`
	Passwd         string `json:"passwd" validate:"required" label:"密码"`
	Port           int    `json:"port" validate:"required,min=1,max=65535" label:"端口"`
}
