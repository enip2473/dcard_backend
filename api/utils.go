package api

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// parseUintQueryParam parses a uint query parameter from the request.
func ParseUintQueryParam(c *fiber.Ctx, name string, defaultValue uint) (uint, error) {
	strValue := c.Query(name)
	if strValue == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseUint(strValue, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid value for query parameter '%s': %v", name, err)
	}
	return uint(value), nil
}
