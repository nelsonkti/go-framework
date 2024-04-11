package databese

import (
	"go-framework/util/xsql/config"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Engine struct {
	Gorm  map[string]*gorm.DB
	Mongo map[string]*mongo.Client
}

type DatabaseClient interface {
	Name() string
	Connect(c map[string]config.DBConfig)
	ConnType(database string) bool
	Result(c *Engine)
}
