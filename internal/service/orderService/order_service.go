package orderService

import "go-framework/internal/dao/order_dao"

type OrderService struct {
	dao *order_dao.OrderDao
}

func NewOrderService(dao *order_dao.OrderDao) *OrderService {
	return &OrderService{dao: dao}
}

func (s *OrderService) Get() {
	s.dao.QueryAll()
}
