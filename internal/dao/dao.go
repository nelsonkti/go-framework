package dao

import (
	"github.com/google/wire"
	"go-framework/internal/dao/order"
	"go-framework/internal/dao/user"
)

var Provider = wire.NewSet(user.NewUserDao, order.NewOrderDao)
