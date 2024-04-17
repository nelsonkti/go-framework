package container

import (
	"go-framework/internal/dao/order_dao"
	"go-framework/internal/service/orderService"
)

type Container struct {
	OrderService *orderService.OrderService
}

func Register() *Container {
	return &Container{OrderService: orderService.NewOrderService(order_dao.NewOrderDao())}
}
