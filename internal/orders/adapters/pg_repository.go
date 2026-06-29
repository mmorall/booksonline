package adapters

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/mmorall/booksonline/internal/orders"

	"github.com/google/uuid"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, o *orders.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	orderQuery := `INSERT INTO orders (id, customer_email, total_price, status, created_at) 
	               VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, orderQuery, o.ID, o.CustomerEmail, o.TotalPrice, o.Status, o.CreatedAt)
	if err != nil {
		return err
	}

	itemQuery := `INSERT INTO order_items (id, order_id, product_id, quantity, price, generated_asset) 
	              VALUES ($1, $2, $3, $4, $5, $6)`
	for _, item := range o.Items {
		_, err = tx.ExecContext(ctx, itemQuery, item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price, item.GeneratedAsset)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	orderQuery := `SELECT id, customer_email, total_price, status, created_at FROM orders WHERE id = $1`
	var o orders.Order
	err := r.db.QueryRowContext(ctx, orderQuery, id).Scan(&o.ID, &o.CustomerEmail, &o.TotalPrice, &o.Status, &o.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, orders.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	itemQuery := `SELECT id, order_id, product_id, quantity, price, generated_asset FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, itemQuery, o.ID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	for rows.Next() {
		var item orders.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.GeneratedAsset); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]*orders.Order, error) {
	query := `SELECT id, customer_email, total_price, status, created_at FROM orders ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var orderList []*orders.Order
	for rows.Next() {
		var o orders.Order
		if err := rows.Scan(&o.ID, &o.CustomerEmail, &o.TotalPrice, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orderList = append(orderList, &o)
	}

	return orderList, nil
}
