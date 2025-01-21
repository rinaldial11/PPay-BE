package controllers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/dto"
	"github.com/ppay/internal/initializers"
	"github.com/ppay/internal/models"
	"github.com/ppay/lib"
)

// Auth godoc
// @Schemes
// @Description Registrasi Account
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Input Email"
// @Param password formData string true "Input Password"
// @Success 201 {object} dto.RegisterDTO
// @Router /auth/register [post]
func Register(c *gin.Context) {
	response := lib.NewResponse(c)
	var input dto.RegisterDTO

	if err := c.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}
	fmt.Println(input)
	if !IsValidEmail(input.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !IsValidPassword(input.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	hasher := lib.GenerateHash(input.Password)
	input.Password = hasher

	// Mulai transaksi
	tx := initializers.DB.Begin()

	// Buat User
	var user = models.User{
		Email:    strings.ToLower(input.Email),
		Password: input.Password,
	}

	// Periksa apakah email sudah digunakan
	if _, err := GetUserByEmail(input.Email); err == nil {
		response.BadRequest("Email is already registered", nil)
		return
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create user", nil)
		return
	}

	// Buat Wallet untuk User
	var wallet = models.Wallet{
		UserID:  user.ID,
		Balance: 0.00, // Saldo awal
	}

	if err := tx.Create(&wallet).Error; err != nil {
		tx.Rollback()
		response.InternalServerError("Failed to create user", nil)
		return
	}

	// Commit transaksi
	tx.Commit()

	response.Created("Success register user", nil)
}

// Auth godoc
// @Schemes
// @Description Login Account
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Input Email"
// @Param password formData string true "Input Password"
// @Success 200 {object} dto.LoginDTO
// @Router /auth/login [post]
func Login(c *gin.Context) {
	response := lib.NewResponse(c)
	var input dto.LoginDTO

	if err := c.ShouldBind(&input); err != nil {
		fmt.Println(err)
	}
	user, err := GetUserByEmail(input.Email)
	if err != nil {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	//compare password
	if user == nil || !lib.VerifyHash(input.Password, user.Password) {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	// Generate token JWT
	token, err := lib.GenerateToken(int(user.ID))
	if err != nil {
		response.InternalServerError("Failed to generate token", nil)
		return
	}

	response.Success("Login successful", gin.H{"token": token})
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Query ke database untuk mencari user berdasarkan email
	if err := initializers.DB.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error; err != nil {
		fmt.Println("error query", err.Error())
		return nil, err // Kembalikan error jika user tidak ditemukan
	}
	return &user, nil
}

func GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	// Query ke database untuk mencari user berdasarkan email
	if err := initializers.DB.Where("phone = ? AND is_deleted = ?", phone, false).First(&user).Error; err != nil {
		return nil, err // Kembalikan error jika user tidak ditemukan
	}
	return &user, nil
}

// Auth godoc
// @Schemes
// @Description Authentication Pin
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Security ApiKeyAuth
// @Param pin formData string true "Input Pin"
// @Success 201 {object} dto.PinDTO
// @Router /auth/pin [post]
func VerifPin(c *gin.Context) {
	response := lib.NewResponse(c)
	// Get user ID from context
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

	var input dto.PinDTO
	if err := c.ShouldBind(&input); err != nil {
		fmt.Println(err)
	}

	user, err := GetUserByIDParam(id)
	if err != nil {
		response.BadRequest("User not found", nil)
		return
	}
	if user.Pin == nil || !lib.VerifyHash(input.Pin, *user.Pin) {
		response.BadRequest("Invalid PIN", nil)
		return
	}
	response.Success("Pin Valid", nil)
}

func IsValidEmail(email string) bool {
	if len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func CheckPassword(c *gin.Context) {
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

	// Cari data user berdasarkan ID
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}
	fmt.Println("Existing User:", user)

	// Bind input data
	var req dto.ExistingPasswordDTO
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	//compare password
	if !lib.VerifyHash(*req.Password, user.Password) {
		response.BadRequest("Invalid password", nil)
		return
	}

	response.Success("Correct password", nil)
}

func CheckAvailPin(c *gin.Context) {
	response := lib.NewResponse(c)

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
	fmt.Println("Existing User:", user)

	// Cek user mempunyai pin atau tidak
	var userPin dto.AvailablePinDTO
	if err := initializers.DB.Model(&models.User{}).Select("pin").Where("id = ? AND is_deleted = ?", user.ID, false).First(&userPin).Error; err != nil {
		response.NotFound("User does not have pin", nil)
		return
	}

	// log.Println(userPin)
	if userPin == (dto.AvailablePinDTO{}) {
		response.NotFound("User does not have pin", nil)
		return
	}

	response.Success("User has pin", nil)
}
