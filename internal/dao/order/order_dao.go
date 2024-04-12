package order

import "fmt"

type OrderDao struct {
}

func NewOrderDao() *OrderDao {
	return &OrderDao{}
}

func (d *OrderDao) QueryAll() {
	fmt.Println("OrderDao QueryAll")
}
