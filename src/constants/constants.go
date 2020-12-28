package constants

const (
	GPC_VERSION = "0.0.1"
	GPS_VERSION = "0.0.1"
)

const (
	WS_EVENT_INIT = "init"
	WS_EVENT_DATA = "data"
	WS_EVENT_ERR  = "err"
	TIME_TEMPLATE = "2006-01-02 15:04:05"
)

const (
	CLIENT_SHELL_TYPE = iota
)

const (
	SYSTEM_MAC     = "darwin"
	SYSTEM_LINUX   = "linux"
	SYSTEM_WINDOWS = "windows"
)

const (
	RUNTIME_PATH        = "/runtime/"
	CONFIG_PATH         = "/config/"
	GPC_CONFIG_FILENAME = "gpc.yaml"
	GPS_CONFIG_FILENAME = "gps.yaml"
	GPC_PID_PATH        = RUNTIME_PATH + "gpc_pid/"
	GPS_PID_PATH        = RUNTIME_PATH + "gps_pid/"
	PID_FILENAME        = "pid"
)

const DEFAULT_SUBPACKAGE = 1000 // 默认分包大小，字节
