package internal

import (
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/reader"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/repository"
)

type Server struct {
	ProductRepo *repository.ProductRepository
	OrderRepo   *repository.OrderRepository
	FileReader  *reader.FileReader
}
