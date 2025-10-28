package controllers

import (
	"net/http"
	"strconv"

	"github.com/dekanayake/kart-challenge/backend-challenge/internal/config"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/repository"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductRepo repository.ProductRepository
}

func NewProductController(productRepo repository.ProductRepository) *ProductController {
	return &ProductController{
		ProductRepo: productRepo,
	}
}

func (p *ProductController) ListProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		config.Logger.Info().
			Str("page", c.Query("page")).
			Msg("Invalid page parameter, defaulting to 1")
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		config.Logger.Info().
			Str("limit", c.Query("limit")).
			Msg("Invalid limit parameter, defaulting to 5")
		limit = 5
	}

	repo := repository.GetProductRepository()
	pageResult, repoErr := repo.ListProducts(page, limit)
	if repoErr != nil {
		config.Logger.Error().
			Err(repoErr).
			Int("page", page).
			Int("limit", limit).
			Msg("Repository error while listing products")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to retrieve products"})
		return
	}

	products := make([]Product, 0)
	for _, product := range pageResult.Items {
		products = append(products, Product{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price.Price,
			Category: product.Category.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":     pageResult.Page,
		"limit":    pageResult.Limit,
		"total":    pageResult.Total,
		"products": products,
	})
}

func (p *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("productId")

	if id == "" {
		config.Logger.Warn().Msg("Missing product ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing product ID"})
		return
	}

	repo := repository.GetProductRepository()
	product, err := repo.GetProductByID(id)
	if err != nil {
		config.Logger.Error().
			Err(err).
			Str("productId", id).
			Msg("Repository error while fetching product")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch product"})
		return
	}

	if product == nil {
		config.Logger.Warn().Str("productId", id).Msg("Product not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, Product{
		ID:       product.ID,
		Name:     product.Name,
		Price:    product.Price.Price,
		Category: product.Category.Name,
	})
}
