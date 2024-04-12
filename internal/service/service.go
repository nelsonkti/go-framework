package service

import (
	"github.com/google/wire"
	"go-framework/internal/service/order"
	"go-framework/internal/service/user"
)

var Provider = wire.NewSet(user.NewUserService, order.NewOrderService)
