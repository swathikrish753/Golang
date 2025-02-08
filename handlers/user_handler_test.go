package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"GolangWorld/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SignUp(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserHandler_SignUp(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("successful signup", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password",
		}
		userJSON, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("SignUp", mock.Anything, &user).Return(nil)

		if assert.NoError(t, handler.SignUp(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.JSONEq(t, `{"message":"User created successfully"}`, rec.Body.String())
		}

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.SignUp(c)
		if assert.Error(t, err) {
			httpError, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusBadRequest, httpError.Code)
				assert.Equal(t, "Invalid request payload", httpError.Message)
			}
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		user := models.User{

			Email:    "swathi",
			Password: "password",
		}
		userJSON, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("SignUp", mock.Anything, &user).Return(assert.AnError)

		err := handler.SignUp(c)
		if assert.Error(t, err) {
			httpError, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusInternalServerError, httpError.Code)
				assert.Equal(t, assert.AnError.Error(), httpError.Message)
			}
		}

		mockService.AssertExpectations(t)
	})
}
