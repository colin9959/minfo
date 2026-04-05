//go:build linux

// Package system 提供 Linux 平台的进程组控制逻辑。

package system

import (
	"os/exec"
	"syscall"
	"time"
)

// setCommandProcessGroup 会为子进程创建独立进程组，方便超时或取消时一起终止。
func setCommandProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

// killCommandProcessGroup 会优先按进程组终止子进程，避免遗留孤儿进程。
func killCommandProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pid := cmd.Process.Pid
	if pid <= 0 {
		return
	}

	// Negative PID targets the whole process group.
	_ = syscall.Kill(-pid, syscall.SIGTERM)
	time.Sleep(300 * time.Millisecond)
	_ = syscall.Kill(-pid, syscall.SIGKILL)
}
