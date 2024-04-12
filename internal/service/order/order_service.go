package order

import (
	"go-framework/internal/dao/order"
)

type OrderService struct {
	dao *order.OrderDao
}

func NewOrderService(dao *order.OrderDao) *OrderService {
	return &OrderService{dao: dao}
}

func (s *OrderService) Get() {
	s.dao.QueryAll()
}
