package entity

import "time"

type Wallet struct {
	Id        int32     `json:"id"`
	UserId    int32     `json:"user_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Balance   float32   `json:"balance" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
