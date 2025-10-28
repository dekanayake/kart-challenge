package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dekanayake/kart-challenge/backend-challenge/internal/config"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/reader"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/repository"
)

type OrderController struct {
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
	FileReader  reader.FileReader
}

func NewOrderController(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, reader reader.FileReader) *OrderController {
	return &OrderController{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		FileReader:  reader,
	}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req OrderReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		config.Logger.Warn().Err(err).Msg("invalid order payload")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if len(req.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order must contain at least one item"})
		return
	}

	if req.CouponCode != "" {
		length := len(req.CouponCode)
		if length < 8 || length > 10 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Coupon code is invalid"})
			return
		}

		valid, err := c.FileReader.SearchPromo(context.Background(), req.CouponCode)
		if err != nil {
			config.Logger.Err(err).Msg("Coupon code validation failed")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Coupon code validation failed"})
			return
		}
		config.Logger.Debug().Bool("valid", valid).Msg("Coupon code validated")
		if !valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Coupon code is invalid"})
			return
		}
	}

	orderItems := make([]repository.OrderItem, 0)

	for _, orderItem := range req.Items {
		product, err := c.ProductRepo.GetProductByID(orderItem.ProductID)
		if err != nil {
			config.Logger.Err(err).Str("product id", orderItem.ProductID).Msg("Error occured while retrieving product for id")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create order"})
		}
		if product == nil {
			config.Logger.Error().Str("product id", orderItem.ProductID).Msg("Product does not exist for the order item")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "product not found in order item"})
		}

		orderItems = append(orderItems, repository.OrderItem{
			ProductId: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		})
	}

	createdOrder, err := c.OrderRepo.CreateOrder(orderItems, req.CouponCode)
	if err != nil {
		config.Logger.Error().Err(err).Msg("failed to create order")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create order"})
		return
	}

	config.Logger.Info().
		Str("orderId", createdOrder.ID).
		Str("couponCode", createdOrder.CouponCode).
		Int("items", len(req.Items)).
		Msg("Order created successfully")

	itemsResponse := make([]OrderItem, 0)
	productsResponse := make([]Product, 0)

	for _, item := range createdOrder.Items {
		itemsResponse = append(itemsResponse, OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
		product, err := c.ProductRepo.GetProductByID(item.ProductId)
		if err != nil {
			config.Logger.Err(err).Str("product id", item.ProductId).Msg("Error occured while retrieving product for id")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating order"})
		}
		productsResponse = append(productsResponse, Product{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price.Price,
			Category: product.Category.Name,
		})
	}

	orderResponse := Order{
		ID:       createdOrder.ID,
		Items:    itemsResponse,
		Products: productsResponse,
	}

	ctx.JSON(http.StatusCreated, orderResponse)
}
