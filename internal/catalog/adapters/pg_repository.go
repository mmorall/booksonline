package adapters

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/mmorall/booksonline/internal/catalog"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*catalog.Product, error) {
	query := `SELECT id, sku, name, price_in_cents, type, stock, is_active 
	          FROM products WHERE id = $1`

	var p catalog.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.PriceInCents, &p.Type, &p.Stock, &p.IsActive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, catalog.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]*catalog.Product, error) {
	query := `SELECT id, sku, name, price_in_cents, type, stock, is_active FROM products`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var products []*catalog.Product
	for rows.Next() {
		var p catalog.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.PriceInCents, &p.Type, &p.Stock, &p.IsActive); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *PostgresRepository) UpdateStock(ctx context.Context, id uuid.UUID, newStock int) error {
	query := `UPDATE products SET stock = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, newStock, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return catalog.ErrProductNotFound
	}
	return nil
}
