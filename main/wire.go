//go:build wireinject
// +build wireinject

package main

import (
	"go-framework/internal/dao"
	"go-framework/internal/service"
)

func InitializeApp() (*App, error) {
	wire.Build(NewApp, service.Provider, dao.Provider)
	return nil, nil
}
