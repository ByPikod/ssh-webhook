package main

import (
	"github.com/ByPikod/ssh-webhook/core"
	"github.com/ByPikod/ssh-webhook/ssh"
)

func main() {
	config := core.InitializeConfig()
	sshCenter := ssh.SSHClientCenter{}
	sshCenter.Initialize(config.YamlConfig.SSH)
}
