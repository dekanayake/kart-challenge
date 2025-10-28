package repository

import (
	"math"
)

type ProductRepository interface {
	GetProductByID(id string) (*Product, error)
	ListProducts(page, limit int) (PaginatedResult[Product], error)
}

type InMemoryProductRepository struct {
	products []Product
}

func newInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: []Product{
			{
				ID:   "1",
				Name: "Waffle with Berries",
				Price: ProductPrice{
					Price: 6.5,
				},
				Category: Category{Name: "Waffle"},
			},
			{
				ID:   "2",
				Name: "Vanilla Bean Crème Brûlée",
				Price: ProductPrice{
					Price: 7,
				},
				Category: Category{Name: "Crème Brûlée"},
			},
			{
				ID:   "3",
				Name: "Macaron Mix of Five",
				Price: ProductPrice{
					Price: 8,
				},
				Category: Category{Name: "Macaron"},
			},
			{
				ID:   "4",
				Name: "Classic Tiramisu",
				Price: ProductPrice{
					Price: 5.5,
				},
				Category: Category{Name: "Tiramisu"},
			},
			{
				ID:   "5",
				Name: "Pistachio Baklava",
				Price: ProductPrice{
					Price: 4,
				},
				Category: Category{Name: "Baklava"},
			},
			{
				ID:   "6",
				Name: "Lemon Meringue Pie",
				Price: ProductPrice{
					Price: 5,
				},
				Category: Category{Name: "Pie"},
			},
			{
				ID:   "7",
				Name: "Red Velvet Cake",
				Price: ProductPrice{
					Price: 4.5,
				},
				Category: Category{Name: "Cake"},
			},
			{
				ID:   "8",
				Name: "Salted Caramel Brownie",
				Price: ProductPrice{
					Price: 4.5,
				},
				Category: Category{Name: "Brownie"},
			},
			{
				ID:   "9",
				Name: "Vanilla Panna Cotta",
				Price: ProductPrice{
					Price: 6.5,
				},
				Category: Category{Name: "Panna Cotta"},
			},
		},
	}
}

func (r *InMemoryProductRepository) GetProductByID(id string) (*Product, error) {
	for _, p := range r.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, nil
}

// TODO: In a DB  limit with offset will be sufficient for small number of products ex : 100k
// but if its in millions will need to handle the pagination using cursor . i.e., return the last cursor of paginated records,
// limit will do a full page scan for every request which can slow down when fetching for pages at tail
// so subsequant pages can fetch from that cursor upto limit.
func (r *InMemoryProductRepository) ListProducts(page, limit int) (PaginatedResult[Product], error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}

	total := len(r.products)
	start := (page - 1) * limit
	end := int(math.Min(float64(start+limit), float64(total)))

	if start >= total {
		return PaginatedResult[Product]{
			Page:  page,
			Limit: limit,
			Total: total,
			Items: []Product{},
		}, nil
	}

	return PaginatedResult[Product]{
		Page:  page,
		Limit: limit,
		Total: total,
		Items: r.products[start:end],
	}, nil
}
