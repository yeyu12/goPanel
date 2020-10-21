package validations

type MachineAdd struct {
	Flag int `json:"flag" validate:"required" label:"类型"`
}

type MachineAddDir struct {
	Name string `json:"name" validate:"required" label:"目录名称"`
}

type MachineAddComputer struct {
	MachineGroupId int64  `json:"machine_group_id" validate:"min=0" label:"目录ID"`
	Alias          string `json:"alias" validate:"required" label:"名称"`
	Host           string `json:"host" validate:"required" label:"地址"`
	User           string `json:"user" validate:"required" label:"用户名"`
	Port           int    `json:"port" validate:"required,min=1,max=65535" label:"端口"`
}
