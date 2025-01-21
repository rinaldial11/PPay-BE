package dto

import "time"

type UserSummaryDTO struct {
	Id       int     `json:"id"`
	Image    *string `json:"image"`
	Email    string  `json:"email,omitempty"`
	Fullname string  `json:"fullname"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Fullname  *string   `json:"fullname" form:"fullname"`
	Email     *string   `json:"email" form:"email" binding:"omitempty,email"`
	Password  *string   `json:"-" form:"password" binding:"omitempty,min=6"`
	Pin       *string   `json:"pin" form:"pin" binding:"omitempty,min=6,max=6"`
	Phone     *string   `json:"phone" form:"phone"`
	Image     *string   `json:"image"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
}

type CreatUserDTO struct {
	Fullname *string `json:"fullname" form:"fullname"`
	Email    string  `json:"email" form:"email" binding:"required,email"`
	Password string  `json:"password" form:"password" binding:"required,min=6"`
	Pin      *string `json:"pin" form:"pin" binding:"omitempty,min=6,max=6"`
	Phone    *string `json:"phone" form:"phone" binding:"required"`
	Image    *string `json:"image"`
}

type ExistingPasswordDTO struct {
	Password *string `json:"-" form:"exist_password" db:"password"`
}

type AvailablePinDTO struct {
	Pin string `json:"pin" db:"pin"`
}
