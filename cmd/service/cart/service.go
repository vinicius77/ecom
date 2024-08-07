package cart

import (
	"fmt"

	"github.com/vinicius77/ecom/types"
)

func getCartItemIDS(items []types.CartItem) ([]int, error) {

	productsIDS := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		productsIDS[i] = item.ProductID
	}

	return productsIDS, nil
}

/*
  - Check if all products are in stock
  - Calculate the total price
	- Reduce quantity of the products in the db
	- Create the order
	- Create the order items

	TODO.: Refactor that later int oa DB transation */
func (h *Handler) CreateOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {

	productsMap := make(map[int]types.Product)

	for _, product := range products {
		productsMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(items, productsMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, productsMap)

	// Refactor the snippet below into a join tables (multiple requests) e.g. orders_items table
	for _, item := range items {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity

		if err := h.productStore.UpdateProduct(product); err != nil {
			return 0, 0, err
		}
	}

	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Dummy address",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil

}

func checkIfCartIsInStock(cartItems []types.CartItem, productsMap map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, cartItem := range cartItems {
		product, ok := productsMap[cartItem.ProductID]

		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", cartItem.ProductID)
		}

		if product.Quantity < cartItem.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, productsMap map[int]types.Product) float64 {
	var total float64

	for _, cartItem := range cartItems {
		product := productsMap[cartItem.ProductID]
		total += float64(cartItem.Quantity) * product.Price
	}

	return total
}
