package handlers

import (
	"net/http"

	"GolangWorld/models"
	"GolangWorld/services"

	"github.com/labstack/echo/v4"
)

// UserHandler struct follows SRP: It only handles HTTP requests.
type UserHandler struct {
	service services.UserService
}

// NewUserHandler initializes a new UserHandler
// DIP: Depends on abstraction (UserService), not a concrete implementation.
func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// SignUp handles user registration requests
// SRP: Handles HTTP request, validation, and response formatting.
func (h *UserHandler) SignUp(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Using request context for lifecycle management (Ensures request doesn't leak)
	err := h.service.SignUp(c.Request().Context(), &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User created successfully"})
}

func (h *UserHandler) Login(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	token, err := h.service.Login(c.Request().Context(), &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
