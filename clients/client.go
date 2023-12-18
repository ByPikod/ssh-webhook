package clients

import (
	"fmt"
	"net"
	"time"

	"github.com/ByPikod/ssh-webhook/core"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

type (
	// The SSH server struct
	SSHClient struct {
		// Configurations
		data     core.SSHServerData
		auth     goph.Auth
		authType string
		callback ssh.HostKeyCallback
		// Keep alive
		retried int
		// The SSH sshClient
		sshClient     *ssh.Client
		highSSHClient *goph.Client
	}
)

func (client *SSHClient) Initialize(data core.SSHServerData, callback ssh.HostKeyCallback) {
	// Set the data
	client.data = data
	client.callback = callback
	// Default values
	client.defaultSettings()
	// Configure the authentication method
	client.configureAuth()
	// Connect to the SSH server
	if *data.KeepAlive {
		log.Infof("Keep alive enabled, connecting to \"%s\"...", data.Name)
		client.keepAlive()
	}
}

// Set the default values for the SSH server
func (client *SSHClient) defaultSettings() {
	// Default values
	if client.data.Retry == nil {
		client.data.Retry = new(bool)
		*client.data.Retry = true
	}
	if client.data.RetryCount == nil {
		client.data.RetryCount = new(int)
		*client.data.RetryCount = 3
	}
	if client.data.RetryInterval == nil {
		client.data.RetryInterval = new(int)
		*client.data.RetryInterval = 5
	}
	if client.data.Timeout == nil {
		client.data.Timeout = new(int)
		*client.data.Timeout = 20
	}
	if client.data.KeepAlive == nil {
		client.data.KeepAlive = new(bool)
		*client.data.KeepAlive = true
	}
	if client.data.Port == nil {
		client.data.Port = new(uint)
		*client.data.Port = 22
	}
	client.retried = 0
}

// Configure the authentication data and method
func (client *SSHClient) configureAuth() error {
	if client.data.PrivateKey != "" {
		// Authenticate with private key and passphrase
		var err error
		client.auth, err = goph.Key(client.data.PrivateKey, client.data.Passphrase)
		if err != nil {
			return err
		}
		client.authType = "key"
	} else {
		// Authenticate with password
		client.auth = goph.Password(client.data.Password)
		client.authType = "password"
	}

	// Authentication method successfully configured
	return nil
}

func (client *SSHClient) retry() {
	if !*client.data.Retry {
		log.Warnf("Retrying is disabled for \"%s\", connection will not be retried", client.data.Name)
		return
	}
	if client.retried >= *client.data.RetryCount {
		log.Warnf("Failed to connect to \"%s\" after %d attempts, connection will not be retried", client.data.Name, client.retried)
		return
	}
	client.retried += 1
	log.Infof("Attempt #%d: Retrying to connect to \"%s\" in %d seconds...", client.retried, client.data.Name, *client.data.RetryInterval)
	time.Sleep(time.Duration(*client.data.RetryInterval) * time.Second)
	client.keepAlive()
}

func (client *SSHClient) connect() (*goph.Client, *ssh.Client, error) {
	// Create the Low level SSH client
	var bannerCallback ssh.BannerCallback
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(client.data.Host, fmt.Sprint(*client.data.Port)), &ssh.ClientConfig{
		User:            client.data.User,
		Auth:            client.auth,
		Timeout:         time.Duration(*client.data.Timeout) * time.Second,
		HostKeyCallback: client.callback,
		BannerCallback:  bannerCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
		},
	})
	if err != nil {
		log.Errorf("Failed to connect to \"%s\": %s", client.data.Name, err.Error())
		log.Debugf("Type of error: %T", err)
		return nil, nil, err
	}

	highSSHClient := &goph.Client{
		Client: sshClient,
	}

	return highSSHClient, sshClient, err
}

// Connect to the SSH server
// Connect is private because its handled by the SSHClient itself
func (client *SSHClient) keepAlive() {
	var err error
	// Create the known hosts callback
	callback, err := goph.DefaultKnownHosts()
	if err != nil {
		log.Errorf("Failed to check known hosts for \"%s\": %s", client.data.Name, err.Error())
		return
	}

	// Create the Low level SSH client
	var bannerCallback ssh.BannerCallback
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(client.data.Host, fmt.Sprint(*client.data.Port)), &ssh.ClientConfig{
		User:            client.data.User,
		Auth:            client.auth,
		Timeout:         time.Duration(*client.data.Timeout) * time.Second,
		HostKeyCallback: callback,
		BannerCallback:  bannerCallback,
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
		},
	})
	if err != nil {
		log.Errorf("Failed to connect to \"%s\": %s", client.data.Name, err.Error())
		log.Debugf("Type of error: %T", err)
		client.retry()
		return
	}

	// Create the high level SSH client
	highSSHClient := &goph.Client{
		Client: sshClient,
	}

	client.sshClient = sshClient
	client.highSSHClient = highSSHClient

	log.Infof("Connection established for server: \"%s\"", client.data.Name)
}

// Execute a command on the SSH server
func (client *SSHClient) Execute(command string) {

}
