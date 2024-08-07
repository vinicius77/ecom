package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vinicius77/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)

	for rows.Next() {
		product, err := ScanRowsIntoProducts(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil

}

func ScanRowsIntoProducts(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt)

	if err != nil {
		return nil, err
	}

	return product, nil

}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec(
		"INSERT INTO products (name, description, image, price, quantity) VALUES (?,?,?,?,?)",
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		product.Quantity)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProductByID(productIDS []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDS)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// Convert ProductIDS into []interface{ }
	args := make([]interface{}, len(productIDS))
	for i, v := range productIDS {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := ScanRowsIntoProducts(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("Update products SET name = ?, price = ?, image = ?, description = ?, price = ?, quantity = ? WHERE id = ?",
		product.Name,
		product.Price,
		product.Image,
		product.Description,
		product.Price,
		product.Quantity,
		product.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
