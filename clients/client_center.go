package clients

import (
	"sync"

	"github.com/ByPikod/ssh-webhook/core"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melbahja/goph"
)

type (
	// The SSH server struct
	SSHClientCenter struct {
		Servers       []*SSHClient
		Configuration core.SSHConfiguration
	}
)

// Initialize the SSH client center
func (s *SSHClientCenter) Initialize(configuration core.SSHConfiguration) {
	s.Configuration = configuration
	log.Info("SSH Client Center initialized.")
	s.createConnections()
	log.Info("SSH connections are created.")
}

// Create the SSH connections
func (s *SSHClientCenter) createConnections() {
	// Create the known hosts callback
	callback, err := goph.DefaultKnownHosts()
	if err != nil {
		log.Errorf("Failed to check known hosts: %s", err.Error())
	}
	// Create a wait group
	var wg sync.WaitGroup
	// Initialize the client in a goroutine
	for _, server := range s.Configuration.Servers {
		wg.Add(1)
		go func(server core.SSHServerData) {
			defer wg.Done()
			client := &SSHClient{}
			client.Initialize(server, callback)
			s.Servers = append(s.Servers, client)
		}(server)
	}
	wg.Wait()
}
