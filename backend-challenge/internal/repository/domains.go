package repository

import "time"

type PaginatedResult[T any] struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Items []T `json:"items"`
}

type Category struct {
	Name string `json:"category" example:"Waffle"`
}

type ProductPrice struct {
	Price float64 `json:"price" example:"12.50"`
}

type Product struct {
	ID       string       `json:"id" example:"10"`
	Name     string       `json:"name" example:"Chicken Waffle"`
	Price    ProductPrice `json:"price" binding:"required"`
	Category Category     `json:"category" binding:"required"`
}

type OrderItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type Order struct {
	ID         string      `json:"id" example:"0000-0000-0000-0000"`
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
	CreatedAt  time.Time   `json:"createdAt"`
}
