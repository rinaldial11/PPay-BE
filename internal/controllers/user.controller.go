package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ppay/internal/dto"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

// @Summary Add User
// @Description This API endpoint is used to create a new user and their wallet.
// @Tags Users
// @Accept json
// @Produce json
// @Param fullname formData string true "Full Name"
// @Param email formData string true "Email Address"
// @Param password formData string true "Password"
// @Param pin formData string true "PIN (6-digit)"
// @Param phone formData string true "Phone Number"
// @Success 200 {object} dto.UserSummaryDTO
// @Router /users [post]

// Create User and Wallet
func CreateUser(c *gin.Context) {
	response := lib.NewResponse(c)
	var input dto.CreatUserDTO
	file, _ := c.FormFile("image")

	// Validate input
	if err := c.ShouldBind(&input); err != nil {
		// Split the error string into individual validation errors
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			response.BadRequest("Invalid input", err.Error())
			return
		}

		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "Email":
				if fieldError.Tag() == "required" {
					response.BadRequest("Email is required", nil)
					return
				} else if fieldError.Tag() == "email" {
					response.BadRequest("Invalid email format", nil)
					return
				}
			case "Password":
				if fieldError.Tag() == "required" {
					response.BadRequest("Password is required", nil)
					return
				} else if fieldError.Tag() == "min" {
					response.BadRequest("Password must be at least 6 characters long", nil)
					return
				}
			case "Phone":
				if fieldError.Tag() == "required" {
					response.BadRequest("Phone is required", nil)
					return
				}
			case "Pin":
				if fieldError.Tag() == "min" {
					response.BadRequest("Pin must be at least 6 characters long", nil)
					return
				} else if fieldError.Tag() == "max" {
					response.BadRequest("Pin must be no more than 6 characters long", nil)
					return
				}
			default:
				response.BadRequest("Invalid input", fieldError.Error())
				return
			}
		}
	}

	if input.Password != "" {
		hasher := lib.GenerateHash(input.Password)
		input.Password = hasher
	}
	if input.Pin != nil {
		hashPin := lib.GenerateHash(*input.Pin)
		input.Pin = &hashPin
	}

	if _, err := GetUserByEmail(input.Email); err == nil {
		response.BadRequest("Email is already registered", nil)
		return
	}

	if _, err := GetUserByPhone(*input.Phone); err == nil {
		response.BadRequest("Phone is already registered", nil)
		return
	}

	if file != nil {
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		maxSize := int64(2 << 20) // 2 MB
		uploadDir := "public/images"

		imagePath, err := lib.UploadImage(c, file, allowedExts, maxSize, uploadDir)
		if err != nil {
			response.BadRequest("Failed to upload image", err.Error())
			return
		}

		input.Image = &imagePath
	} else {
		imageDefault := ""
		input.Image = &imageDefault
	}

	// Mulai transaksi
	tx := initializers.DB.Begin()

	// Buat User
	user := models.User{
		Fullname: *input.Fullname,
		Email:    input.Email,
		Password: input.Password,
		Pin:      input.Pin,
		Phone:    input.Phone,
		Image:    input.Image,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Buat Wallet untuk User
	wallet := models.Wallet{
		UserID:  user.ID,
		Balance: 0.00, // Saldo awal
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	// Commit transaksi
	tx.Commit()

	response.Created("Success create user", dto.UserSummaryDTO{
		Id:       int(user.ID),
		Email:    user.Email,
		Fullname: user.Fullname,
		Phone:    user.Phone,
		Image:    user.Image,
	})

}

// Users godoc
// @Summary Users
// @Description  Get All Users
// @Tags Users
// @Accept json
// @Produce json
// @Param search query string false "Search Users"
// @Param page query int false "Page Users"
// @Param limit query int false "Limit Users"
// @Param sort query string false "Sort Users"
// @Param order query string false "Order Users"
// @Success 200 {object} lib.Response{data=[]dto.UserSummaryDTO,pageInfo=lib.PageInfo}
// @Router /users [get]
// Get All Users
func GetUsers(c *gin.Context) {
	response := lib.NewResponse(c)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	search := c.DefaultQuery("search", "")
	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")

	// Validate input for page and limit
	if page < 1 {
		response.BadRequest("Invalid input", "Page must be 1 or greater")
		return
	}
	if limit < 1 {
		response.BadRequest("Invalid input", "Limit must be 1 or greater")
		return
	}

	// Validate sorting parameters
	if sort != "id" && sort != "fullname" && sort != "phone" {
		response.BadRequest("Invalid input", "Sort must be 'id', 'fullname', or 'phone'")
		return
	}
	if order != "asc" && order != "desc" {
		response.BadRequest("Invalid input", "Order must be 'asc' or 'desc'")
		return
	}

	offset := (page - 1) * limit

	// Query database with filters and pagination
	var users []dto.UserSummaryDTO
	query := initializers.DB.Model(&models.User{}).
		Select("id, fullname, phone, image").
		Where("is_deleted = ?", false)

	// Apply search filter if provided
	if search != "" {
		query = query.Where("fullname ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply sorting, limit, and offset
	if err := query.Order(sort + " " + order).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		response.InternalServerError("Failed to get users", nil)
		return
	}

	// Count total users matching the criteria
	var totalCount int64
	countQuery := initializers.DB.Model(&models.User{}).Where("is_deleted = ?", false)
	if search != "" {
		countQuery = countQuery.Where("fullname ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	countQuery.Count(&totalCount)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Build pagination info
	pageInfo := &lib.PageInfo{
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		TotalPage:   totalPages,
		TotalData:   int(totalCount),
	}
	if page >= totalPages {
		pageInfo.NextPage = 0
	}
	if page <= 1 {
		pageInfo.PrevPage = 0
	}

	// Return response
	response.GetAllSuccess("Success get user", users, pageInfo)
}

// Users godoc
// @Schemes
// @Description  Get Profile
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserSummaryDTO
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	response := lib.NewResponse(c)

	// Get user ID from context
	userId, exists := c.Get("UserId")
	fmt.Println(userId)
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	log.Println(userId)
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Use DTO for selected fields
	var userSummary dto.UserSummaryDTO

	// Query only required fields
	if err := initializers.DB.Model(&models.User{}).
		Select("id, email, image, fullname, phone").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&userSummary).Error; err != nil {
		response.NotFound("User not found", nil)
		return
	}

	// Return response with filtered data
	response.Success("Success get user", userSummary)
}

// Users godoc
// @Schemes
// @Description Update Movies
// @Tags Users
// @Accept mpfd
// @Produce json
// @Security ApiKeyAuth
// @Param fullname formData string false "Update Full Name"
// @Param email formData string false "Update Email"
// @Param password formData string false "Update Password"
// @Param pin formData string false "Update Pin"
// @Param phone formData string false "Update Phone"
// @Param image formData file false "Update Image"
// @Success 200 {object} lib.Response{data=dto.UpdateUserRequest}
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	response := lib.NewResponse(c)
	file, _ := c.FormFile("image")

	// Ambil userId dari konteks
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Cari data user berdasarkan ID
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	// Bind input data
	var req dto.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			response.BadRequest("Invalid input", err.Error())
			return
		}
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "Phone":
				if fieldError.Tag() == "registered" {
					response.BadRequest("Phone is registered", nil)
					return
				}
			case "Pin":
				if fieldError.Tag() == "min" {
					response.BadRequest("Pin must be at least 6 characters long", nil)
					return
				} else if fieldError.Tag() == "max" {
					response.BadRequest("Pin must be no more than 6 characters long", nil)
					return
				}
			default:
				response.BadRequest("Invalid input", fieldError.Error())
				return
			}
		}
	}

	// Update data hanya jika ada input
	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		if !IsValidPassword(*req.Password) {
			response.BadRequest("Password must be at least 8 characters long", nil)
			return
		}
		hashedPassword := lib.GenerateHash(*req.Password)
		user.Password = hashedPassword
	}
	if req.Pin != nil {
		hashPin := lib.GenerateHash(*req.Pin)
		user.Pin = &hashPin
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if file != nil {
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		maxSize := int64(2 << 20) // 2MB
		uploadDir := "public/images"

		// if *user.Image != "" {
		// 	oldFilePath := *user.Image
		// 	if err := os.Remove(oldFilePath); err != nil {
		// 		log.Printf("Failed to delete old profile picture: %s", err)
		// 	}
		// }

		imagePath, err := lib.UploadImage(c, file, allowedExts, maxSize, uploadDir)
		if err != nil {
			response.BadRequest("Failed to upload image", err.Error())
			return
		}
		user.Image = &imagePath
	}

	// Perbarui waktu
	user.UpdatedAt = time.Now()

	// Simpan perubahan ke database
	if err := initializers.DB.Save(&user).Error; err != nil {
		response.InternalServerError("Failed to update user", err.Error())
		return
	}

	// Respon sukses
	response.Success("Update user success", dto.UserSummaryDTO{
		Id:       int(user.ID),
		Email:    user.Email,
		Fullname: user.Fullname,
		Image:    user.Image,
		Phone:    user.Phone,
	})
}

func GetUserByIDParam(userID int) (*models.User, error) {
	var user models.User
	if err := initializers.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetTargetUserById(c *gin.Context) {
	response := lib.NewResponse(c)
	id, _ := strconv.Atoi(c.Param("id"))

	detailsUser, _ := GetUserByIDParam(id)

	response.Success("Details user", detailsUser)
}

// @Delete User godoc
// @Summary User
// @Description Delete User
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /users/{id} [delete]
// Delete User
func DeleteUser(c *gin.Context) {
	response := lib.NewResponse(c)

	// Ambil userId dari konteks
	userId, exists := c.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	// Cari data user
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	// Hapus user (soft delete)
	user.IsDeleted = true
	user.UpdatedAt = time.Now()
	if err := initializers.DB.Save(&user).Error; err != nil {
		response.InternalServerError("Failed to delete user", err.Error())
		return
	}

	// Respon sukses
	response.Success(fmt.Sprintf("User with ID %d deleted successfully", id), nil)
}
