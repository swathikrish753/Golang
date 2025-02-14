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

func (m *MockUserService) Login(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
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
func TestUserHandler_Login(t *testing.T) {
	e := echo.New()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("successful login", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password",
		}
		userJSON, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Login", mock.Anything, &user).Return("mockToken", nil)

		if assert.NoError(t, handler.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"token":"mockToken"}`, rec.Body.String())
		}

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Login(c)
		if assert.Error(t, err) {
			httpError, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusBadRequest, httpError.Code)
				assert.Equal(t, "Invalid request payload", httpError.Message)
			}
		}
	})

	t.Run("unauthorized error", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		userJSON, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Login", mock.Anything, &user).Return("", assert.AnError)

		err := handler.Login(c)
		if assert.Error(t, err) {
			httpError, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusUnauthorized, httpError.Code)
				assert.Equal(t, assert.AnError.Error(), httpError.Message)
			}
		}

		mockService.AssertExpectations(t)
	})
}
