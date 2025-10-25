package handler

import (
	"echto/internal/model"
	"echto/internal/service"
	"echto/pkg/logger"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

// GetUsers handles GET /api/v1/users
// @Summary Get all users
// @Description Retrieve a paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} model.UserListResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Get users from service
	users, err := h.userService.GetUsers(page, limit)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get users")
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_server_error",
			Message: "Failed to get users",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, users)
}

// GetUser handles GET /api/v1/users/:id
// @Summary Get user by ID
// @Description Retrieve a specific user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.UserResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	// Parse user ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	// Get user from service
	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
		}
		logger.Log.Error().Err(err).Msg("Failed to get user")
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_server_error",
			Message: "Failed to get user",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /api/v1/users
// @Summary Create a new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.UserCreateRequest true "User data"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req model.UserCreateRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	// Create user
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusConflict, model.ErrorResponse{
				Error:   "email_exists",
				Message: "Email already exists",
				Code:    http.StatusConflict,
			})
		}
		logger.Log.Error().Err(err).Msg("Failed to create user")
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_server_error",
			Message: "Failed to create user",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser handles PUT /api/v1/users/:id
// @Summary Update user
// @Description Update an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.UserUpdateRequest true "User data"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	// Parse user ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	var req model.UserUpdateRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	// Update user
	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
		}
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusConflict, model.ErrorResponse{
				Error:   "email_exists",
				Message: "Email already exists",
				Code:    http.StatusConflict,
			})
		}
		logger.Log.Error().Err(err).Msg("Failed to update user")
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_server_error",
			Message: "Failed to update user",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /api/v1/users/:id
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	// Parse user ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
	}

	// Delete user
	if err := h.userService.DeleteUser(uint(id)); err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
		}
		logger.Log.Error().Err(err).Msg("Failed to delete user")
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "internal_server_error",
			Message: "Failed to delete user",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.NoContent(http.StatusNoContent)
}
