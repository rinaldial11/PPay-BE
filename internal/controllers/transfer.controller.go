package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

	type ResponseTRF struct {
		Amount     float64
		TargetUser int
		Type       string
	}

// Transaction godoc
// @Summary Transfer
// @Description  transfer money
// @Tags Transaction
// @Accept x-www-form-urlencoded
// @Produce json
// @Param id path string false "send to"
// @Param amount formData string false "amount of money"
// @Success 200 {object} ResponseTRF
// @Security ApiKeyAuth
// @Router /transfer/{id} [post]
func Transfer(c *gin.Context) {
	var input struct {
		Amount float64 `form:"amount" json:"amount" binding:"required"`
	}

	response := lib.NewResponse(c)
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized access", nil)
		return
	}

	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Validate input
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// Begin a database transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin database transaction"})
		return
	}

	// Fetch the sender's wallet to check the balance
	var senderWallet models.Wallet
	if err := tx.Where("user_id = ?", id).First(&senderWallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to retrieve sender's wallet", err.Error())
		return
	}

	// Check if sender has enough balance for the transfer
	if senderWallet.Balance < input.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Create the main transaction record
	transaction := models.Transaction{
		UserID:          uint(id),
		Amount:          input.Amount,
		TransactionType: "transfer",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction", "details": err.Error()})
		return
	}

	targetId, _ := strconv.Atoi(c.Param("id"))

	// Fetch the target user to ensure valid recipient
	var targetUser models.User
	if err := initializers.DB.Where("id = ? AND is_deleted = false", targetId).First(&targetUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
		return
	}

	// Create the transfer transaction record
	transfer := models.TransferTransaction{
		TransactionID: transaction.ID,
		TargetUserID:  uint(targetId),
	}

	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create transfer transaction", err.Error())
		return
	}

	// Update sender's wallet balance
	senderWallet.Balance -= input.Amount
	if err := tx.Save(&senderWallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to update sender's wallet balance", err.Error())
		return
	}

	// Update target user's wallet balance
	var targetWallet models.Wallet
	if err := tx.Where("user_id = ?", targetId).First(&targetWallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to retrieve target user's wallet", err.Error())
		return
	}
	targetWallet.Balance += input.Amount
	if err := tx.Save(&targetWallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to update target user's wallet balance", err.Error())
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to commit database transaction", err.Error())
		return
	}

	response.Created("Success transfer", ResponseTRF{
		Amount:     input.Amount,
		TargetUser: targetId,
		Type:       transaction.TransactionType,
	})
}
