package orderService

import (
	"go-framework/internal/dao/order_dao"
	"go-framework/util/xerror"
)

type OrderService struct {
	dao *order_dao.OrderDao
}

func NewOrderService(dao *order_dao.OrderDao) *OrderService {
	return &OrderService{dao: dao}
}

func (s *OrderService) Get() error {
	s.dao.QueryAll()
	return xerror.BadRequest(3002, "没有参数")
}
