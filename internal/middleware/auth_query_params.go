package middleware

import "github.com/gofiber/fiber/v2"

func AuthQueryParams(c *fiber.Ctx) error {
	// get query params
	userId := c.Query("user_id")

	// If the cookie is not set, return an unauthorized status
	if userId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	// If the token is valid, pass the request to the next handler
	return c.Next()
}
