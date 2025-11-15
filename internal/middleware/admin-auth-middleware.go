package middleware

import (
	"app/internal/common"
	"app/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminAuth(userService *services.UserService, adminService *services.AdminService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			firebaseUID, ok := c.Get(common.AUTH_USER_ID).(string)
			if !ok || firebaseUID == "" {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "user identity not found in context",
				})
			}

			user, err := userService.GetUserProfile(c.Request().Context(), firebaseUID)
			if err != nil {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "user profile not found",
				})
			}

			isAdmin, err := adminService.IsAdmin(c.Request().Context(), user.ID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "could not verify admin status",
				})
			}

			if !isAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "admin access required",
				})
			}

			return next(c)
		}
	}
}

func RequireAdminRole(userService *services.UserService, adminService *services.AdminService) echo.MiddlewareFunc {
	return AdminAuth(userService, adminService)
}
