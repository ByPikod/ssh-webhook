package core

import (
	"io"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/yaml.v3"
)

type (
	// The SSH configuration struct
	SSHConfiguration struct {
		Servers []SSHServerData `yaml:"servers"`
	}

	// The SSH server struct
	SSHServerData struct {
		Name       string `yaml:"name"`
		Host       string `yaml:"host"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		PrivateKey string `yaml:"path_to_key"`
		Passphrase string `yaml:"passphrase"`
		// Optional
		Port          *uint `yaml:"port"`
		KeepAlive     *bool `yaml:"keep_alive"`
		Retry         *bool `yaml:"retry"`
		RetryCount    *int  `yaml:"retry_count"`
		RetryInterval *int  `yaml:"retry_interval"`
		Timeout       *int  `yaml:"timeout"`
	}

	// The auth struct for the webhook
	WebhookAuthData struct {
		IPWhitelist  []string          `yaml:"ip_whitelist"`
		IPBlacklist  []string          `yaml:"ip_blacklist"`
		IPFromHeader string            `yaml:"ip_from_header"`
		Headers      map[string]string `yaml:"headers"`
	}

	// The webhook struct
	WebhookData struct {
		Path     string          `yaml:"path"`
		SSHName  string          `yaml:"ssh_name"`
		WhiteCMD []string        `yaml:"cmd_whitelist"`
		BlackCMD []string        `yaml:"cmd_blacklist"`
		Auth     WebhookAuthData `yaml:"auth"`
	}

	// The YAML config struct
	YamlConfig struct {
		SSH      SSHConfiguration `yaml:"ssh"`
		Webhooks []WebhookData    `yaml:"webhooks"`
	}

	// The config struct
	Config struct {
		Port       string     // The port to listen on
		Host       string     // The host to listen on
		Yaml       string     // The path to the YAML config file
		YamlConfig YamlConfig // The YAML config
	}
)

// Return the first non-empty string
func or(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

// Initialize the config
func InitializeConfig() *Config {
	// Create the config instance
	config := &Config{
		Port: or(os.Getenv("SSH_WEBHOOK_PORT"), "3000"),
		Host: or(os.Getenv("SSH_WEBHOOK_HOST"), ""),
		Yaml: or(os.Getenv("SSH_WEBHOOK_YAML"), "ssh-webhook.yaml"),
	}

	// Load the YAML config
	yamlFile, err := os.Open(config.Yaml)
	if err != nil {
		log.Errorf("Error opening YAML config file: %s", err)
		os.Exit(0)
	}

	defer yamlFile.Close()
	yamlByte, err := io.ReadAll(yamlFile)
	if err != nil {
		log.Errorf("Error reading YAML config file: %s", err)
		os.Exit(0)
	}

	yml, err := ParseConfig(yamlByte)
	if err != nil {
		log.Errorf("Error parsing YAML config file: %s", err)
		os.Exit(0)
	}
	config.YamlConfig = *yml
	return config
}

func ParseConfig(yamlByte []byte) (*YamlConfig, error) {
	// Default config
	yamlConfig := &YamlConfig{}
	// Parse the YAML config
	err := yaml.Unmarshal(yamlByte, yamlConfig)
	if err != nil {
		return nil, err
	}
	return yamlConfig, nil
}
