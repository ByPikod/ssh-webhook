package main

import (
	"github.com/ByPikod/ssh-webhook/clients"
	"github.com/ByPikod/ssh-webhook/core"
	"github.com/ByPikod/ssh-webhook/router"
)

func main() {
	config := core.InitializeConfig()
	sshCenter := clients.SSHClientCenter{}
	// Initialize all the SSH clients asynchronously
	sshCenter.Initialize(config.YamlConfig.SSH)
	// Initialize the HTTP server
	router.Listen(config, &sshCenter)
}
