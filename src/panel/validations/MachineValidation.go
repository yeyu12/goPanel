package validations

type MachineAdd struct {
	Flag           int    `json:"flag" validate:"required" label:"类型"`
	MachineGroupId int64  `json:"machine_group_id" validate:"required,gte=0" label:"目录ID"`
	Alias          string `json:"alias" validate:"required" label:"名称"`
	Host           string `json:"host" validate:"required"`
	User           string `json:"user" validate:"required" label:"用户名"`
	Port           int    `json:"port" validate:"required,gt=0,lte=65535" label:"端口"`
	Name           string `json:"name" validate:"required" label:"目录名称"`
}
