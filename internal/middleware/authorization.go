package middleware

import "github.com/gofiber/fiber/v2"

func AuthRole(authorizedRoles []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
				"success": false,
			})
		}

		for _, authorizedRole := range authorizedRoles {
			if role == authorizedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"success": false,
		})
	}
}
