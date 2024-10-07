package entity

import "time"

type TopUpRequest struct {
	Id     int     `json:"id" name:"id"`
	Amount float64 `json:"amount" name:"amount"`
	Notes  string  `json:"notes" name:"notes"`
}

type TransactionResponse struct {
	Message string `json:"message"`
}

type TransactionRequest struct {
	SourceId      int     `json:"source_id" name:"source_id"`
	DestinationId int     `json:"destination_id" name:"destination_id"`
	Type          string  `json:"type" name:"type"`
	Category      string  `json:"category" name:"category"`
	Amount        float64 `json:"amount" name:"amount"`
	Notes         string  `json:"notes" name:"notes"`
}

type TransferRequest struct {
	SourceId      int     `json:"source_id" name:"source_id"`
	DestinationId int     `json:"destination_id" name:"destination_id"`
	Amount        float64 `json:"amount" name:"amount"`
	Notes         string  `json:"notes" name:"notes"`
}

type TransactionGetResponse struct {
	Id        int       `json:"Id" name:"Id"`
	UserId    int       `json:"user_id" name:"user_id"`
	Amount    float64   `json:"amount" name:"amount"`
	Category  string    `json:"category" name:"category"`
	Type      string    `json:"type" name:"type"`
	Notes     string    `json:"notes" name:"notes"`
	CreatedAt time.Time `json:"createdAt" name:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" name:"updatedAt"`
}

type TransactionGetRequest struct {
	Type     string `json:"type" name:"type"`
	Category string `json:"category" name:"category"`
	UserId   int    `json:"user_id" name:"user_id"`
	Size     int    `json:"size" name:"size"`
	Page     int    `json:"page" name:"page"`
}

type Pagination struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	PageSize  int `json:"page_size"`
	Page      int `json:"page"`
}

type TransactionGetResponseWithPagination struct {
	Data       []TransactionGetResponse `json:"data"`
	Pagination Pagination               `json:"pagination"`
}
