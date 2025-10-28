package controllers

type Product struct {
	ID       string  `json:"id" example:"10"`
	Name     string  `json:"name" example:"Chicken Waffle"`
	Price    float64 `json:"price" example:"12.50"`
	Category string  `json:"category" example:"Waffle"`
}

type OrderItem struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type OrderReq struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items" binding:"required"`
}

type Order struct {
	ID       string      `json:"id" example:"0000-0000-0000-0000"`
	Items    []OrderItem `json:"items"`
	Products []Product   `json:"products"`
}
