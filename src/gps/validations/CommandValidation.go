package validations

type CommandAdd struct {
	Command      string   `json:"command" validate:"required" label:"命令"`
	Flag         string   `json:"flag" validate:"required" label:"执行方式"`
	IsType       int      `json:"is_type" validate:"required"`
	PlanExecTime string   `json:"plan_exec_time" label:"计划执行时间"`
	Ids          []string `json:"ids" label:"ID"`
}
