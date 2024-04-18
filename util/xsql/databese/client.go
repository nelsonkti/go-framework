package databese

import (
	"context"
	"fmt"
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
	Close()
}

func (e *Engine) Close() {
	e.gormClose()
	e.mongodbClose()
}

func (e *Engine) gormClose() {
	for _, g := range e.Gorm {
		db, err := g.DB()
		if err != nil || db.Ping() != nil {
			fmt.Println(err)
			continue
		}
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (e *Engine) mongodbClose() {
	for _, m := range e.Mongo {
		err := m.Disconnect(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
