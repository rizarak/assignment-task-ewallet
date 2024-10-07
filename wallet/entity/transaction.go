package entity

import "time"

type Transaction struct {
	Id            int       `json:"id"`
	SourceId      int       `json:"source_id" gorm:"not null" binding:"required"`
	DestinationId int       `json:"destination_id" gorm:"not null" binding:"required"`
	Amount        float64   `json:"amount" gorm:"not null" binding:"required"`
	Type          string    `gorm:"not null" binding:"required,oneof=debt credit"`
	Category      string    `gorm:"not null" binding:"required,oneof=topup transfer withdrawal"`
	Notes         string    `json:"notes"`
	Timestamp     time.Time `gorm:"autoCreateTime" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TransactionRequest struct {
	SourceId      int     `json:"source_id" name:"source_id"`
	DestinationId int     `json:"destination_id" name:"destination_id"`
	Type          string  `json:"type" name:"type"`
	Category      string  `json:"category" name:"category"`
	Amount        float64 `json:"amount" name:"amount"`
	Notes         string  `json:"notes"`
}

type TransactionGetRequest struct {
	Type     string `json:"type" name:"type"`
	Category string `json:"category" name:"category"`
	UserId   int    `json:"user_id" name:"user_id"`
	Size     int    `json:"size" name:"size"`
	Page     int    `json:"page" name:"page"`
}

type TransactionResponse struct {
	Transaction []Transaction `json:"transaction"`
	Count       int64         `json:"count"`
}

type Pagination struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	PageSize  int `json:"page_size"`
	Page      int `json:"page"`
}

type TransactionGetResponseWithPagination struct {
	Data       []Transaction `json:"data"`
	Pagination Pagination    `json:"pagination"`
}
