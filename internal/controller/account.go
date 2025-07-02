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

type createAccountRequest struct {
	Currency    string `json:"currency"`
	PhoneNumber string `json:"phone_number" binding:required,max="12"`
}

// createAccountHandler godoc
// @Summary Create a new account
// @Description Creates a new account for the authenticated user with the provided currency and phone number
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body createAccountRequest true "Account data"
// @Security BearerAuth
// @Success 201 {object} models.Account
// @Failure 400 {object} map[string]string "Invalid input, user not found, account already exists, or invalid currency"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts [post]
func createAccountHandler(ctx *gin.Context) {
	const op = "createAccountHandler"

	var req createAccountRequest

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
		logger.Error.Printf("%s: userID conversion error: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	account := &models.Account{
		UserID:      userID,
		Currency:    req.Currency,
		PhoneNumber: req.PhoneNumber,
	}

	err = service.CreateAccount(account)
	if err != nil {
		logger.Error.Printf("%s: service.CreateAccount: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "such user does not exist, firstly create user"})
		case errors.Is(err, errs.ErrAccountAlreadyExists):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account with such credentials already exits"})
		case errors.Is(err, errs.ErrInvalidCurrency):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency was sent"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"account": account})
}

// updateAccountHandler godoc
// @Summary Update an existing account
// @Description Updates an account's phone number, balance, or currency for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body models.UpdateAccount true "Account update data"
// @Security BearerAuth
// @Success 200 {object} models.Account
// @Failure 400 {object} map[string]string "Invalid input, account not found, or invalid currency"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts [patch]
func updateAccountHandler(ctx *gin.Context) {
	const op = "updateAccountHandler"

	var req models.UpdateAccount

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
		logger.Error.Printf("%s: userID conversion error: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to convert userID to int"})
		return
	}

	account := &models.UpdateAccount{
		ID:          req.ID,
		PhoneNumber: req.PhoneNumber,
		Balance:     req.Balance,
		Currency:    req.Currency,
		UserID:      userID,
	}

	newAccount, err := service.UpdateAccount(account)
	if err != nil {
		logger.Error.Printf("%s: service.UpdateAccount: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		case errors.Is(err, errs.ErrInvalidCurrency):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency was sent"})
		case errors.Is(err, errs.ErrFraud):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot change others account"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"account": newAccount})
}

type deleteAccountRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// deleteAccountHandler godoc
// @Summary Delete an existing account
// @Description Deletes an account for the authenticated user by account ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Security BearerAuth
// @Success 202 {object} map[string]string "Account deleted successfully"
// @Failure 400 {object} map[string]string "Invalid input, account not found, or non-zero balance"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/{id} [delete]
func deleteAccountHandler(ctx *gin.Context) {
	const op = "deleteAccountHandler"

	var req deleteAccountRequest

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

	err = service.DeleteAccount(req.ID, userID)
	if err != nil {
		logger.Error.Printf("%s: service.DeleteAccount: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		case errors.Is(err, errs.ErrNoZeroBalance):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account's balance is not zero"})
		case errors.Is(err, errs.ErrFraud):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete other's data"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "account deleted successfully"})
}

type getAccountByIDRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// getAccountByIDHandler godoc
// @Summary Get an account by ID
// @Description Retrieves an account by its ID for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Security BearerAuth
// @Success 200 {object} models.Account
// @Failure 400 {object} map[string]string "Invalid input or account not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/{id} [get]
func getAccountByIDHandler(ctx *gin.Context) {
	const op = "getAccountByIDHandler"

	var req getAccountByIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindUri: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	user, err := service.GetAccountByID(req.ID)
	if err != nil {
		logger.Error.Printf("%s: service.GetAccountByID: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"account": user})
}

type getAccountByUserIDRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// getAccountByUserIDHandler godoc
// @Summary Get an account by user ID
// @Description Retrieves an account by its ID for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} models.Account
// @Failure 400 {object} map[string]string "Invalid input or account not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/users/{id} [get]
func getAccountByUserIDHandler(ctx *gin.Context) {
	const op = "getAccountByUserIDHandler"

	var req getAccountByUserIDRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindUri: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	user, err := service.GetAccountByUserID(req.ID)
	if err != nil {
		logger.Error.Printf("%s: service.GetAccountByID: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account with such user id not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"account": user})
}

// getInActiveAccountsHandler godoc
// @Summary Get inactive accounts
// @Description Retrieves a list of inactive accounts
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Account
// @Failure 400 {object} map[string]string "Invalid input or unauthorized"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/inactive [get]

func getInActiveAccountsHandler(ctx *gin.Context) {
	const op = "getInActiveAccountsHandler"

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

	users, err := service.GetInactiveAccounts()
	if err != nil {
		logger.Error.Printf("%s: service.GetInactiveUsers: %v", op, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"accounts": users})
}

type getAccountByCurrencyRequest struct {
	Currency string `json:"currency" binding:"required"`
}

// getAccountByCurrency godoc
// @Summary Get accounts by currency
// @Description Retrieves all accounts with the specified currency for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param currency body getAccountByCurrencyRequest true "Currency"
// @Security BearerAuth
// @Success 200 {array} models.Account
// @Failure 400 {object} map[string]string "Invalid input or invalid currency"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/currency [get]
func getAccountByCurrency(ctx *gin.Context) {
	const op = "getAccountByCurrency"

	var req getAccountByCurrencyRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindUri: %v", op, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read sent data"})
		return
	}

	accounts, err := service.GetAccountsByCurrency(req.Currency)
	if err != nil {
		logger.Error.Printf("%s: service.GetAccountsByCurrency: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrInvalidCurrency):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency was sent"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

type getAccountBalanceRequest struct {
	ID int `uri:"id" binding:"required,min=1"`
}

// getAccountBalanceHandler godoc
// @Summary Get an account balance by ID
// @Description Retrieves an account balance by its ID for the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Security BearerAuth
// @Success 200 {object} models.Account
// @Failure 400 {object} map[string]string "Invalid input or account not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /accounts/{id}/balance [get]
func getAccountBalanceHandler(ctx *gin.Context) {
	const op = "getAccountBalanceHandler"

	var req getAccountBalanceRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		logger.Error.Printf("%s: ctx.ShouldBindUri: %v", op, err)
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

	balance, err := service.GetAccountBalance(req.ID, userID, userID == service.AdminID)
	if err != nil {
		logger.Error.Printf("%s: service.GetAccountBalance: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		case errors.Is(err, errs.ErrFraud):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot access others account balance"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}
