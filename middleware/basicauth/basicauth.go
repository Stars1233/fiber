package basicauth

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// The key for the username value stored in the context
const (
	usernameKey contextKey = iota
)

const basicScheme = "Basic"

// New creates a new middleware handler
func New(config Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	// Return new handler
	return func(c fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get authorization header and ensure it matches the Basic scheme
		auth := utils.Trim(c.Get(fiber.HeaderAuthorization), ' ')
		if auth == "" || len(auth) > cfg.HeaderLimit {
			return cfg.Unauthorized(c)
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || !utils.EqualFold(parts[0], basicScheme) {
			return cfg.Unauthorized(c)
		}

		// Decode the header contents
		raw, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return cfg.Unauthorized(c)
		}

		// Get the credentials
		var creds string
		if c.App().Config().Immutable {
			creds = string(raw)
		} else {
			creds = utils.UnsafeString(raw)
		}

		// Check if the credentials are in the correct form
		// which is "username:password".
		index := strings.Index(creds, ":")
		if index == -1 {
			return cfg.Unauthorized(c)
		}

		// Get the username and password
		username := creds[:index]
		password := creds[index+1:]

		if cfg.Authorizer(username, password, c) {
			c.Locals(usernameKey, username)
			return c.Next()
		}

		// Authentication failed
		return cfg.Unauthorized(c)
	}
}

// UsernameFromContext returns the username found in the context
// returns an empty string if the username does not exist
func UsernameFromContext(c fiber.Ctx) string {
	username, ok := c.Locals(usernameKey).(string)
	if !ok {
		return ""
	}
	return username
}
