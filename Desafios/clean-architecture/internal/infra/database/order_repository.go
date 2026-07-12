package database

import (
	"database/sql"

	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/domain/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.DB.Prepare(
		"INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	return err
}

func (r *OrderRepository) GetAll() ([]entity.Order, error) {
	rows, err := r.DB.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.Price, &o.Tax, &o.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
