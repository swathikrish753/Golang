package routes

import (
	"GolangWorld/handlers"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes sets up API routes
// SRP: Keeps route management separate from handlers and main.
func RegisterRoutes(e *echo.Group, userHandler *handlers.UserHandler) {
	e.POST("/signup", userHandler.SignUp)
}
