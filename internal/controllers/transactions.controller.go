package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

// type TransactionHistoryResponse struct {
// 	ID              uint    `json:"id"`
// 	Amount          float64 `json:"amount"`
// 	TransactionType string  `json:"transaction_type"`
// 	Notes           *string `json:"notes,omitempty"`
// 	CreatedAt       string  `json:"created_at"`
// 	FromUser        *string `json:"from_user,omitempty"`
// 	ToUser          *string `json:"to_user,omitempty"`
// }

type IncomeResponse struct {
	Income string `json:"income"`
}

type ExpenseResponse struct {
	Expense string `json:"expense"`
}

// Struct untuk response transaksi
// Struct for transaction response
type TransactionHistoryResponse struct {
	ID                  uint    `json:"id"`
	TransactionType     string  `json:"transaction_type"`
	UserImage           *string `json:"user_image,omitempty"`
	UserFullname        *string `json:"user_fullname,omitempty"`
	UserPhone           *string `json:"user_phone,omitempty"`
	Amount              float64 `json:"amount"`
	RelatedUserImage    *string `json:"related_user_image,omitempty"`
	RelatedUserFullname *string `json:"related_user_fullname,omitempty"`
	RelatedUserPhone    *string `json:"related_user_phone,omitempty"`
	CreatedAt           string  `json:"created_at"`
}

// Transaction godoc
// @Summary get history transaction
// @Schemes
// @Description  Get Transaction History
// @Tags Transaction
// @Produce json
// @Success 200 {object} TransactionHistoryResponse
// @Security ApiKeyAuth
// @Router /transaction/history [get]
// Function untuk mendapatkan riwayat transaksi
// Function to get transaction history
func GetTransactionHistory(c *gin.Context) {
	response := lib.NewResponse(c)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	if page < 1 {
		response.BadRequest("Invalid page number", nil)
		return
	}
	if limit < 1 {
		response.BadRequest("Invalid limit number", nil)
		return
	}

	offset := (page - 1) * limit

	// Get user ID from context
	userID, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query to fetch transactions
	var transactions []TransactionHistoryResponse
	query := `
-- Transfers where the user is the sender (Sent)
SELECT 
    t.id,
    'Sent' AS transaction_type,
    u.image AS user_image,
    u.fullname AS user_fullname,
    u.phone AS user_phone,
    t.amount,
    target_u.image AS related_user_image,
    target_u.fullname AS related_user_fullname,
    target_u.phone AS related_user_phone,
    to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
FROM 
    users u
JOIN 
    transactions t ON u.id = t.user_id
JOIN 
    transfer_transactions tt ON tt.transaction_id = t.id
JOIN 
    users target_u ON target_u.id = tt.target_user_id
WHERE 
    t.user_id = ?  
    AND (
        u.fullname ILIKE ? OR 
        u.phone ILIKE ? OR 
        target_u.fullname ILIKE ? OR 
        target_u.phone ILIKE ?
    )

UNION ALL

-- Transfers where the user is the recipient (Received)
SELECT 
    t.id,
    'Received' AS transaction_type,
    target_u.image AS user_image,
    target_u.fullname AS user_fullname,
    target_u.phone AS user_phone,
    t.amount,
    u.image AS related_user_image,
    u.fullname AS related_user_fullname,
    u.phone AS related_user_phone,
    to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
FROM 
    users u
JOIN 
    transactions t ON u.id = t.user_id
JOIN 
    transfer_transactions tt ON tt.transaction_id = t.id
JOIN 
    users target_u ON target_u.id = tt.target_user_id
WHERE 
    tt.target_user_id = ?  
    AND (
        target_u.fullname ILIKE ? OR 
        target_u.phone ILIKE ? OR 
        u.fullname ILIKE ? OR 
        u.phone ILIKE ?
    )

UNION ALL

-- Top-up history (Top-Up)
SELECT 
    t.id,
    'Top-Up' AS transaction_type,
    u.image AS user_image,
    u.fullname AS user_fullname,
    u.phone AS user_phone,
    t.amount,
    NULL AS related_user_image,
    pm.name AS related_user_fullname,
    NULL AS related_user_phone,
    to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
FROM 
    users u
JOIN 
    transactions t ON u.id = t.user_id
JOIN 
    topup_transactions tu ON tu.transaction_id = t.id
JOIN 
    payment_methods pm ON pm.id = tu.payment_method_id
WHERE 
    t.user_id = ?  
    AND (
        pm.name ILIKE ?
    )

ORDER BY 
    created_at DESC
LIMIT ? OFFSET ?;`

	searchQuery := "%" + search + "%"

	// Execute query
	if err := initializers.DB.Raw(query, id, searchQuery, searchQuery, searchQuery, searchQuery, id, searchQuery, searchQuery, searchQuery, searchQuery, id, searchQuery, limit, offset).Scan(&transactions).Error; err != nil {
		response.InternalServerError("Failed to retrieve transaction history", err.Error())
		return
	}

	// Query to count total transactions
	var totalCount int64
	countQuery := `
SELECT COUNT(*)
FROM (
    -- Transfers where the user is the sender (Sent)
    SELECT t.id
    FROM users u
    JOIN transactions t ON u.id = t.user_id
    JOIN transfer_transactions tt ON tt.transaction_id = t.id
    JOIN users target_u ON target_u.id = tt.target_user_id
    WHERE t.user_id = ?
      AND (u.fullname ILIKE ? OR u.phone ILIKE ? OR target_u.fullname ILIKE ? OR target_u.phone ILIKE ?)

    UNION ALL

    -- Transfers where the user is the recipient (Received)
    SELECT t.id
    FROM users u
    JOIN transactions t ON u.id = t.user_id
    JOIN transfer_transactions tt ON tt.transaction_id = t.id
    JOIN users target_u ON target_u.id = tt.target_user_id
    WHERE tt.target_user_id = ?
      AND (target_u.fullname ILIKE ? OR target_u.phone ILIKE ? OR u.fullname ILIKE ? OR u.phone ILIKE ?)

    UNION ALL

    -- Top-up history (Top-Up)
    SELECT t.id
    FROM users u
    JOIN transactions t ON u.id = t.user_id
    JOIN topup_transactions tu ON tu.transaction_id = t.id
    JOIN payment_methods pm ON pm.id = tu.payment_method_id
    WHERE t.user_id = ?
      AND pm.name ILIKE ?
) AS total_transactions;`

	// Execute count query
	if err := initializers.DB.Raw(
		countQuery,
		id, searchQuery, searchQuery, searchQuery, searchQuery, // Sent
		id, searchQuery, searchQuery, searchQuery, searchQuery, // Received
		id, searchQuery, // Top-Up
	).Scan(&totalCount).Error; err != nil {
		response.InternalServerError("Failed to count transactions", err.Error())
		return
	}

	// Calculate total pages for pagination
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Create pagination info
	pageInfo := &lib.PageInfo{
		CurrentPage: page,
		NextPage:    0,
		PrevPage:    0,
		TotalPage:   totalPages,
		TotalData:   int(totalCount),
	}
	if page < totalPages {
		pageInfo.NextPage = page + 1
	}
	if page > 1 {
		pageInfo.PrevPage = page - 1
	}

	// Return success response
	response.GetAllSuccess("Success get transaction history", transactions, pageInfo)
}

// GetTransactionHistory retrieves a paginated list of transaction history for a user
// func GetTransactionHistory(c *gin.Context) {
// 	response := lib.NewResponse(c)

// 	// Parse query parameters
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

// 	if page < 1 {
// 		response.BadRequest("Invalid page number", nil)
// 		return
// 	}
// 	if limit < 1 {
// 		response.BadRequest("Invalid limit number", nil)
// 		return
// 	}

// 	offset := (page - 1) * limit

// 	// Get user ID from context
// 	userID, exists := c.Get("UserId")
// 	if !exists {
// 		response.Unauthorized("Unauthorized", nil)
// 		return
// 	}
// 	id, ok := userID.(int)
// 	if !ok {
// 		response.InternalServerError("Failed to parse user ID from token", nil)
// 		return
// 	}

// 	// Query to fetch transaction history with user details
// 	var transactions []TransactionHistoryResponse
// 	query := `
//         SELECT t.id,
//                amount,
//                t.transaction_type,
//                t.notes,
//                to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
//                fu.fullname AS from_user,
//                tu.fullname AS to_user
//         FROM transactions t
//         LEFT JOIN transfer_transactions tt ON t.id = tt.transaction_id
//         LEFT JOIN users fu ON tt.target_user_id = fu.id
//         LEFT JOIN users tu ON t.user_id = tu.id
//         WHERE t.user_id = ? AND t.is_deleted = false
//         ORDER BY t.created_at DESC
//         LIMIT ? OFFSET ?
//     `
// 	if err := initializers.DB.Raw(query, id, limit, offset).Scan(&transactions).Error; err != nil {
// 		response.InternalServerError("Failed to retrieve transaction history", err.Error())
// 		return
// 	}

// 	// Count total transactions for pagination
// 	var totalCount int64
// 	countQuery := `
//         SELECT COUNT(*)
//         FROM transactions t
//         WHERE t.user_id = ? AND t.is_deleted = false
//     `
// 	if err := initializers.DB.Raw(countQuery, id).Scan(&totalCount).Error; err != nil {
// 		response.InternalServerError("Failed to count transactions", err.Error())
// 		return
// 	}

// 	// Calculate total pages
// 	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

// 	// Build pagination info
// 	pageInfo := &lib.PageInfo{
// 		CurrentPage: page,
// 		NextPage:    page + 1,
// 		PrevPage:    page - 1,
// 		TotalPage:   totalPages,
// 		TotalData:   int(totalCount),
// 	}
// 	if page >= totalPages {
// 		pageInfo.NextPage = 0
// 	}
// 	if page <= 1 {
// 		pageInfo.PrevPage = 0
// 	}

// 	// Return success response
// 	response.GetAllSuccess("Success get transaction history", transactions, pageInfo)
// }

// Transaction godoc
// @Summary get user income
// @Schemes
// @Description  get user income
// @Tags Transaction
// @Produce json
// @Success 200 {object} IncomeResponse
// @Security ApiKeyAuth
// @Router /transaction/income [get]
func GetUserIncome(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userID, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query to get income transactions
	var income IncomeResponse
	query := `
        select 
		sum(amount) 
		as income from 
		transactions t where user_id = ? and transaction_type ='top_up'
    `
	if err := initializers.DB.Raw(query, id).Scan(&income).Error; err != nil {
		response.InternalServerError("Failed to retrieve income transactions", err.Error())
		return
	}

	response.Success("Success get user income", income)
}

// Transaction godoc
// @Summary get user expense
// @Schemes
// @Description  get user expense
// @Tags Transaction
// @Produce json
// @Success 200 {object} ExpenseResponse
// @Security ApiKeyAuth
// @Router /transaction/expense [get]
func GetUserExpenses(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userID, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Query to get expense transactions
	var expenses ExpenseResponse
	query := `
        select 
		sum(amount) 
		as expense from 
		transactions t where user_id = ? and transaction_type ='transfer'
    `
	if err := initializers.DB.Raw(query, id).Scan(&expenses).Error; err != nil {
		response.InternalServerError("Failed to retrieve expense transactions", err.Error())
		return
	}

	response.Success("Success get user expenses", expenses)
}

// Transaction godoc
// @Summary get payment method
// @Schemes
// @Description  get payment method
// @Tags Transaction
// @Produce json
// @Success 200 {object} models.PaymentMethod
// @Security ApiKeyAuth
// @Router /transaction/payment [get]
func GetPaymentMethod(c *gin.Context) {
	response := lib.NewResponse(c)

	// Query to get expense transactions
	var paymentMethod []models.PaymentMethod
	query := `
        select 
		id, name, tax 
		from 
		payment_methods
    `
	if err := initializers.DB.Raw(query).Scan(&paymentMethod).Error; err != nil {
		response.InternalServerError("Failed to retrieve payment method", err.Error())
		return
	}

	response.Success("Success get payment method", paymentMethod)
}
