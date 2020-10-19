package validations

type Login struct {
	Username string `validate:"required" json:"username" label:"用户名"`
	Passwd   string `validate:"required" json:"passwd" label:"密码"`
}
