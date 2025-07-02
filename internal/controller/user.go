package controller

import (
	"SB/internal/errs"
	"SB/internal/models"
	"SB/internal/service"
	"SB/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserRequest struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

// createUserHandler godoc
// @Summary Create a new user
// @Description Creates a new user with the provided full name and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body createUserRequest true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Invalid input or user already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/sign-up [post]
func createUserHandler(ctx *gin.Context) {
	const op = "createUserHandler"

	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindJSON: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	user := &models.User{
		FullName: req.FullName,
		Password: req.Password,
	}

	err = service.CreateUser(user)
	if err != nil {
		logger.Error.Printf("%s: service.CreateUser: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrInvalidFullName) || errors.Is(err, errs.ErrPasswordTooShort):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "check sent data for requirements"})
		case errors.Is(err, errs.ErrUserAlreadyExists):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with such credentials already exits"})
		case errors.Is(err, errs.ErrCreateHash):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "password too long"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

// updateUserHandler godoc
// @Summary Update a user
// @Description Updates user information based on the provided ID and optional fields
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UpdateUser true "Updated user data"
// @Security BearerAuth
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid input or user not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [patch]
func updateUserHandler(ctx *gin.Context) {
	const op = "updateUserHandler"
	var req models.UpdateUser

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindJSON: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	userIDAny, ok := ctx.Get(userIDCtx)
	if !ok {
		logger.Error.Printf("%s: userID absent in Context", op)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse userID"})
		return
	}

	userID, ok := userIDAny.(int)
	if !ok {
		logger.Error.Printf("%s: userID conversion error: %v", op, userIDAny)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	if userID != req.ID {
		logger.Error.Printf("%s: someone is trying to change others data, userID token: %d userID request: %d", op, userID, req.ID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	err = service.UpdateUser(&req)
	if err != nil {
		logger.Error.Printf("%s: service.UpdateUser: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrInvalidFullName) || errors.Is(err, errs.ErrPasswordTooShort):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "check sent data for requirements"})
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found for update"})
		case errors.Is(err, errs.ErrCreateHash):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "new password is too long"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "user successfully updated"})
}

type deleteUserRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// deleteUserHandler godoc
// @Summary Delete a user
// @Description Deletes a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid ID or user has dependencies"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func deleteUserHandler(ctx *gin.Context) {
	const op = "deleteUserHandler"

	var req deleteUserRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindJSON: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	userIDAny, ok := ctx.Get(userIDCtx)
	if !ok {
		logger.Error.Printf("%s: userID absent in Context", op)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse userID"})
		return
	}

	userID, ok := userIDAny.(int)
	if !ok {
		logger.Error.Printf("%s: userID conversion error: %v", op, userIDAny)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	if userID != req.ID && userID != service.AdminID {
		logger.Error.Printf("%s: someone is trying to change others data, userID token: %d userID request: %d", op, userID, req.ID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	err = service.DeleteUser(req.ID)
	if err != nil {
		logger.Error.Printf("%s: service.DeleteUser: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrAccountExists):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "delete account of the user first"})
		case errors.Is(err, errs.ErrCreditsExists):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user should pay for credits first"})
		case errors.Is(err, errs.ErrDepositsExists):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "delete deposit first"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "user deleted successfully"})
}

type getUserByIDRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// getUserByIDHandler godoc
// @Summary Get a user by ID
// @Description Retrieves a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid ID or user not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [get]
func getUserByIDHandler(ctx *gin.Context) {
	const op = "getUserByIDHandler"

	var req getUserByIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindUri: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	user, err := service.GetUserByID(req.ID)
	if err != nil {
		logger.Error.Printf("%s: service.GetUserByID: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// getInActiveUsersHandler godoc
// @Summary Get inactive users
// @Description Retrieves a list of inactive users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 400 {object} map[string]string "Invalid input or unauthorized"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/inactive [get]
func getInActiveUsersHandler(ctx *gin.Context) {
	const op = "getInActiveUsersHandler"

	userIDAny, ok := ctx.Get(userIDCtx)
	if !ok {
		logger.Error.Printf("%s: userID absent in Context", op)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse userID"})
		return
	}

	userID, ok := userIDAny.(int)
	if !ok {
		logger.Error.Printf("%s: userID conversion error: %v", op, userIDAny)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	if userID != service.AdminID {
		logger.Error.Printf("%s: someone is trying to get others data, userID token: %d", op, userID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	users, err := service.GetInactiveUsers()
	if err != nil {
		logger.Error.Printf("%s: service.GetInactiveUsers: %v", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

type authenticateRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// authenticateHandler godoc
// @Summary Authenticate a user
// @Description Authenticates a user with full name and password, returns a token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body authenticateRequest true "User credentials"
// @Success 200 {object} map[string]interface{} "Contains user and token"
// @Failure 400 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/sign-in [post]
func authenticateHandler(ctx *gin.Context) {
	const op = "authenticateHandler"

	var req authenticateRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindJSON: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	token, user, err := service.AuthenticateUser(req.FullName, req.Password)
	if err != nil {
		logger.Error.Printf("%s: service.AuthenticateUser: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		case errors.Is(err, errs.ErrInvalidPassword):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password was entered"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

type restoreUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
}

// restoreUserHandler godoc
// @Summary Restore a deleted user
// @Description Restores a user by their full name
// @Tags users
// @Accept json
// @Produce json
// @Param user body restoreUserRequest true "User full name"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid input or user not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/restore [post]
func restoreUserHandler(ctx *gin.Context) {
	const op = "restoreUserHandler"

	var req restoreUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindJSON: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	userIDAny, ok := ctx.Get(userIDCtx)
	if !ok {
		logger.Error.Printf("%s: userID absent in Context", op)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse userID"})
		return
	}

	userID, ok := userIDAny.(int)
	if !ok {
		logger.Error.Printf("%s: userID conversion error: %v", op, userIDAny)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	if userID != service.AdminID {
		logger.Error.Printf("%s: someone is trying to restore others data, userID token: %d", op, userID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	err = service.RestoreUser(req.FullName)
	if err != nil {
		logger.Error.Printf("%s: service.RestoreUser: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user has been restored"})
}

type findUserByNameRequest struct {
	FullName string `json:"full_name" binding:"required"`
}

// findUserByNameHandler godoc
// @Summary Find users by name
// @Description Retrieves users matching the provided full name
// @Tags users
// @Accept json
// @Produce json
// @Param user body findUserByNameRequest true "User full name"
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 400 {object} map[string]string "Invalid input or user not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/find [get]
func findUserByNameHandler(ctx *gin.Context) {
	const op = "findUserByName"

	var req findUserByNameRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindJSON: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	users, err := service.FindUserByName(req.FullName)
	if err != nil {
		logger.Error.Printf("%s: service.FindUserByName: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}
