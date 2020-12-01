package controller

import (
	log "github.com/sirupsen/logrus"
	"goPanel/src/constants"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service"
	"goPanel/src/gpc/service/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"
)

// 连接中继端，和本地ssh连接
func SshConnectRelay(conn *net.TCPConn, message interface{}) {
	log.Info("进入执行器")

	data := message.(service.Message).Data.(map[string]interface{})
	relayClient := ssh.NewRelayClient()
	relayAddr := config.Conf.App.ServerHost + ":" + strconv.Itoa(int(data["port"].(float64)))
	log.Info("连接中继端，relayAddr:", relayAddr)
	err := relayClient.RelayConn(relayAddr, constants.CLIENT_SHELL_TYPE, uint32(data["cols"].(float64)), uint32(data["rows"].(float64)))
	if err != nil {
		log.Error(err)
		return
	}
}

// 设置客户端信息
func SettingClientInfo(conn *net.TCPConn, message interface{}) {
	log.Info("设置客户端信息")
	dataMap := message.(service.Message).Data.(map[string]interface{})
	if dataMap["id"] != config.Conf.App.Uid {
		return
	}

	// 更新配置
	config.Conf.App.LocalName = dataMap["name"].(string)
	config.Conf.Ssh.Username = dataMap["username"].(string)
	config.Conf.Ssh.Password = dataMap["passwd"].(string)
	config.Conf.Ssh.Port = int(dataMap["port"].(float64))

	serverPort, _ := strconv.Atoi(config.Conf.App.ServerPort)

	confApp := map[string]interface{}{
		"app": map[string]interface{}{
			"debug":                   config.Conf.App.Debug,
			"log_level":               config.Conf.App.LogLevel,
			"log_output_type":         config.Conf.App.LogOutputType,
			"log_output_flag":         config.Conf.App.LogOutputFlag,
			"log_path":                config.Conf.App.LogPath,
			"server_host":             config.Conf.App.ServerHost,
			"server_port":             serverPort,
			"local_name":              config.Conf.App.LocalName,
			"control_heartbeat_time":  config.Conf.App.ControlHeartbeatTime,
			"control_reconn_tcp_time": config.Conf.App.ControlReconnTcpTime,
			"uid_path":                config.Conf.App.UidPath,
		},
	}

	confSsh := map[string]interface{}{
		"ssh": map[string]interface{}{
			"username": config.Conf.Ssh.Username,
			"password": config.Conf.Ssh.Password,
			"port":     config.Conf.Ssh.Port,
		},
	}
	// 写入配置文件
	c, err := yaml.Marshal(confApp)
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(config.GpcConfigFileName, c, 0755)
	if err != nil {
		log.Error("配置文件写入错误！")
	}

	c, err = yaml.Marshal(confSsh)
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(config.GpcSshConfigPath, c, 0755)
	if err != nil {
		log.Error("配置文件写入错误！")
	}

	// 重启客户端
}

// 重启客户端主机
func Reboot(conn *net.TCPConn, message interface{}) {
	cmd := exec.Command("reboot")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err)
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		log.Error(err)
	}
}

// 重启服务
func RestartService(conn *net.TCPConn, message interface{}) {

}
