package controller

import (
	"context"
	log "github.com/sirupsen/logrus"
	"goPanel/src/common"
	"goPanel/src/constants"
	"goPanel/src/gpc/config"
	"goPanel/src/gpc/service"
	"goPanel/src/gpc/service/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
)

// 连接中继端，和本地ssh连接
func SshConnectRelay(ctx context.Context, conn *net.TCPConn, message interface{}) {
	data := message.(service.Message).Data.(map[string]interface{})
	relayClient := ssh.NewRelayClient()
	conf := config.NewConf()
	relayAddr := conf.App.ServerHost + ":" + strconv.Itoa(int(data["port"].(float64)))
	log.Info("连接中继端，relayAddr:", relayAddr)
	err := relayClient.RelayConn(ctx, relayAddr, constants.CLIENT_SHELL_TYPE, uint32(data["cols"].(float64)), uint32(data["rows"].(float64)))
	if err != nil {
		log.Error(err)
		return
	}
}

// 设置客户端信息
func SettingClientInfo(ctx context.Context, conn *net.TCPConn, message interface{}) {
	dataMap := message.(service.Message).Data.(map[string]interface{})
	conf := config.NewConf()
	if dataMap["id"] != conf.App.Uid {
		return
	}

	// 更新配置
	conf.App.LocalName = dataMap["name"].(string)
	conf.Ssh.Username = dataMap["username"].(string)
	conf.Ssh.Password = dataMap["passwd"].(string)
	conf.Ssh.Port = int(dataMap["port"].(float64))

	serverPort, _ := strconv.Atoi(conf.App.ServerPort)

	confApp := map[string]interface{}{
		"app": map[string]interface{}{
			"debug":                   conf.App.Debug,
			"log_level":               conf.App.LogLevel,
			"log_output_type":         conf.App.LogOutputType,
			"log_output_flag":         conf.App.LogOutputFlag,
			"log_path":                conf.App.LogPath,
			"server_host":             conf.App.ServerHost,
			"server_port":             serverPort,
			"local_name":              conf.App.LocalName,
			"control_heartbeat_time":  conf.App.ControlHeartbeatTime,
			"control_reconn_tcp_time": conf.App.ControlReconnTcpTime,
			"uid_path":                conf.App.UidPath,
		},
	}

	confSsh := map[string]interface{}{
		"ssh": map[string]interface{}{
			"username": conf.Ssh.Username,
			"password": conf.Ssh.Password,
			"port":     conf.Ssh.Port,
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

	if err := common.SendPidRestart(os.Getpid()); err != nil {
		log.Error(err)
	}
}

// 重启客户端主机
func Reboot(ctx context.Context, conn *net.TCPConn, message interface{}) {
	cmd := exec.Command("sudo", "reboot")
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
func RestartService(ctx context.Context, conn *net.TCPConn, message interface{}) {
	if err := common.SendPidRestart(os.Getpid()); err != nil {
		log.Error(err)
	}
}
