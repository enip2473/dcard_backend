package api

import (
	"dcard_backend/pkg/ads"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateAd(repo ads.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var ad ads.Ad
		if err := c.BodyParser(&ad); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		// Validate request
		if err := ad.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Save ad to database
		if err := repo.CreateAd(&ad); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Send success response
		return c.Status(fiber.StatusCreated).JSON(ad)
	}
}

func ListAds(repo ads.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offset, err := parseUintQueryParam(c, "offset", 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		limit, err := parseUintQueryParam(c, "limit", 10)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		conditions := make(map[string]interface{})
		// Parsing and adding conditions based on query parameters
		if age, err := parseIntQueryParam(c, "age", 0); err == nil && age > 0 {
			conditions["age"] = age
		}
		if gender := c.Query("gender"); gender != "" {
			conditions["gender"] = gender
		}
		// Add other conditions similarly...

		ads, total, err := repo.ListActiveAds(int(offset), int(limit), conditions)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"items": ads,
			"total": total,
		})
	}
}

func parseUintQueryParam(c *fiber.Ctx, name string, defaultValue uint) (uint, error) {
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

func parseIntQueryParam(c *fiber.Ctx, name string, defaultValue int) (int, error) {
	strValue := c.Query(name)
	if strValue == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(strValue, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid value for query parameter '%s': %v", name, err)
	}
	return int(value), nil
}
