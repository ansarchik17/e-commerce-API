package repositories

import (
	"context"
	"e-commerce-API/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductsRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(conn *pgxpool.Pool) *ProductsRepository {
	return &ProductsRepository{db: conn}
}

func (repository *ProductsRepository) Create(ctx context.Context, product models.Product) (int, error) {
	var id int

	err := repository.db.QueryRow(ctx, "insert into products(name, price) values($1, $2) returning id", product.Name, product.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}
