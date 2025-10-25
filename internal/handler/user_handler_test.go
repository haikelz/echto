package handler

import (
	"bytes"
	"echto/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(req *model.UserCreateRequest) (*model.UserResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUser(id uint) (*model.UserResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUsers(page, limit int) (*model.UserListResponse, error) {
	args := m.Called(page, limit)
	return args.Get(0).(*model.UserListResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	args := m.Called(id, req)
	return args.Get(0).(*model.UserResponse), args.Error(1)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserHandler_CreateUser(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		requestBody    model.UserCreateRequest
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name: "successful user creation",
			requestBody: model.UserCreateRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			mockSetup: func(mockService *MockUserService) {
				mockService.On("CreateUser", mock.AnythingOfType("*model.UserCreateRequest")).
					Return(&model.UserResponse{
						ID:    1,
						Name:  "John Doe",
						Email: "john@example.com",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "email already exists",
			requestBody: model.UserCreateRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			mockSetup: func(mockService *MockUserService) {
				mockService.On("CreateUser", mock.AnythingOfType("*model.UserCreateRequest")).
					Return((*model.UserResponse)(nil), errors.New("email already exists"))
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			handler := NewUserHandler(mockService)

			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateUser(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:   "successful user retrieval",
			userID: "1",
			mockSetup: func(mockService *MockUserService) {
				mockService.On("GetUser", uint(1)).
					Return(&model.UserResponse{
						ID:    1,
						Name:  "John Doe",
						Email: "john@example.com",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not found",
			userID: "999",
			mockSetup: func(mockService *MockUserService) {
				mockService.On("GetUser", uint(999)).
					Return((*model.UserResponse)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			err := handler.GetUser(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
		})
	}
}
