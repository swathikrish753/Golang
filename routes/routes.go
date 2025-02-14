package routes

import (
	"GolangWorld/handlers"

	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes sets up API routes
// SRP: Keeps route management separate from handlers and main.
func RegisterRoutes(e *echo.Group, userHandler *handlers.UserHandler) {
	e.POST("/signup", userHandler.SignUp)
}

func RegisterListRoutes(e *echo.Group, listHandler *handlers.ListHandler) {
	e.GET("/getusers", func(c echo.Context) error {
		users, err := listHandler.ListAllUsers(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, users)
	})
}
