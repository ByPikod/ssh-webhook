package router

import (
	"fmt"

	"github.com/ByPikod/ssh-webhook/clients"
	"github.com/ByPikod/ssh-webhook/core"
	"github.com/gofiber/fiber/v2"
)

// Initialize the HTTP server
func Listen(config *core.Config, sshCenter *clients.SSHClientCenter) {
	// Create a new Fiber instance
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Listen on the specified host and port
	app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
}
