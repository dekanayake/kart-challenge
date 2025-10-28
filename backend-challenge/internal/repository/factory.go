package repository

import "sync"

type RepositoryFactory struct {
	productRepo ProductRepository
	orderRepo   OrderRepository
}

var factory RepositoryFactory

var (
	productOnce sync.Once
	orderOnce   sync.Once
)

func GetProductRepository() ProductRepository {
	productOnce.Do(func() {
		factory.productRepo = newInMemoryProductRepository()
	})
	return factory.productRepo
}

func GetOrderRepository() OrderRepository {
	orderOnce.Do(func() {
		factory.orderRepo = newInMemoryOrderRepository()
	})
	return factory.orderRepo
}
