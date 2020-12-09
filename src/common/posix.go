// +build !windows

package common

import "syscall"

// 给pid发送重启信号
func SendPidRestart(pid int) error {
	return syscall.Kill(pid, syscall.SIGUSR2)
}
