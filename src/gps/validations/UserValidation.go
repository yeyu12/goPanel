package validations

type Login struct {
	Username string `validate:"required" json:"username" label:"用户名"`
	Passwd   string `validate:"required" json:"passwd" label:"密码"`
}

type UserAdd struct {
	Username     string `validate:"required" json:"username" label:"用户名"`
	Passwd       string `validate:"required" json:"passwd" label:"密码"`
	RepeatPasswd string `json:"repeat_passwd" label:"重复密码"`
}
