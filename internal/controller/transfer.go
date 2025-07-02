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

type createTransferRequest struct {
	FromAccountID int    `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int    `json:"to_account_id" binding:"required,min=1"`
	Amount        int    `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
}

// createTransferHandler godoc
// @Summary Create a new transfer
// @Description Creates a transfer between two accounts for the authenticated user with the specified amount and currency
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body createTransferRequest true "Transfer data"
// @Security BearerAuth
// @Success 201 {object} models.Transfer
// @Failure 400 {object} map[string]string "Invalid input, invalid user ID, mismatched currencies, or insufficient balance"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transfers [post]
func createTransferHandler(ctx *gin.Context) {
	const op = "createTransferHandler"

	var req createTransferRequest

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

	trnx := &models.Transfer{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	}

	result, err := service.CreateTransfer(trnx, userID)
	if err != nil {
		logger.Error.Printf("%s: service.CreateTransferr: %v", op, err)
		switch {
		case errors.Is(err, errs.ErrNotFound):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "enter valid user_id"})
		case errors.Is(err, errs.ErrInvalidCurrency):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "currencies must match"})
		case errors.Is(err, errs.ErrInsufficientBalance):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": result})
}
