package api

import (
	"dcard_backend/pkg/ads"

	"github.com/gofiber/fiber/v2"
)

func CreateAd(repo ads.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var adInput ads.AdInput
		if err := c.BodyParser(&adInput); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := adInput.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		countries, platforms, ad, err := ParseNewAdInput(adInput)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for _, countryCode := range countries {
			if country, err := repo.FirstOrCreateCountry(countryCode); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			} else {
				ad.Countries = append(ad.Countries, country)
			}
		}

		for _, platformName := range platforms {
			if platform, err := repo.FirstOrCreatePlatform(platformName); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			} else {
				ad.Platforms = append(ad.Platforms, platform)
			}
		}

		if err := repo.CreateAd(&ad); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(ad)
	}
}

func ListAds(repo ads.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query, err := ParseQuery(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		ads, err := repo.ListActiveAds(query)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"items": ads,
		})
	}
}
