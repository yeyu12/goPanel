package validations

type MachineSaveComputer struct {
	Id       string `json:"id" label:"ID"`
	Name     string `json:"name" label:"名称"`
	Username string `json:"username" label:"用户名"`
	Passwd   string `json:"passwd" validate:"required" label:"密码"`
	Port     int    `json:"port" validate:"required,min=1,max=65535" label:"端口"`
}
