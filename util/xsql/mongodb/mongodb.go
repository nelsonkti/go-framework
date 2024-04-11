package mongodb

import (
	"fmt"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/databese"
	"go.mongodb.org/mongo-driver/mongo/options"
)
import "go.mongodb.org/mongo-driver/mongo"
import "context"

type MongoDB struct {
	c      map[string]config.DBConfig
	client map[string]*mongo.Client
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		client: make(map[string]*mongo.Client),
	}
}

func (m *MongoDB) Name() string {
	return "mongodb"
}

func (m *MongoDB) Connect(c map[string]config.DBConfig) {
	m.c = c
	for _, dbConfig := range m.c {
		err := m.connect(dbConfig)
		if err != nil {
			panic(fmt.Sprintf("The database %s connection failed， error: %s", dbConfig.Database, err))
		}
	}
}

func (m *MongoDB) connect(c config.DBConfig) error {
	// 创建客户端选项
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s?authSource=admin",
		c.Username, c.Password, c.Host, c.Port, c.Database,
	))

	// 建立到 MongoDB 的连接
	var client *mongo.Client
	var err error
	if client, err = mongo.Connect(context.Background(), clientOptions); err != nil {
		return err
	}
	// 检查连接是否成功
	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}
	alias := c.Alias
	if alias == "" {
		alias = c.Database
	}

	m.client[alias] = client
	return nil
}

func (m *MongoDB) ConnType(database string) bool {
	if database != "mongodb" {
		return false
	}
	return true
}

func (m *MongoDB) Result(c *databese.Engine) {
	c.Mongo = m.client
}
