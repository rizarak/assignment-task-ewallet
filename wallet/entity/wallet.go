package entity

import "time"

type Wallet struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id" gorm:"not null" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Balance   float64   `json:"balance" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
