package xgorm

import (
	"fmt"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/databese"
	"go-framework/util/xsql/xgorm/clickhouse"
	"go-framework/util/xsql/xgorm/mysql"
	"go-framework/util/xsql/xgorm/postgresql"
	"gorm.io/gorm"
	"time"
)

type DB struct {
	Host   gorm.Dialector
	Master *[]gorm.Dialector
	Slave  *[]gorm.Dialector
}

type Gorm struct {
	DB     map[string]DataBase
	c      map[string]config.DBConfig
	client map[string]*gorm.DB
}

type DataBase interface {
	Conn(config config.DBConfig) (*gorm.DB, error)
}

func NewGorm() *Gorm {
	db := make(map[string]DataBase)
	db["mysql"] = mysql.NewDB()
	db["clickhouse"] = clickhouse.NewDB()
	db["pgsql"] = postgresql.NewDB()
	return &Gorm{
		DB:     db,
		client: make(map[string]*gorm.DB),
	}
}

func (g *Gorm) Name() string {
	return "gorm"
}

func (g *Gorm) Connect(c map[string]config.DBConfig) {
	g.c = c
	databases := make(map[string]*gorm.DB)
	for _, dbConfig := range g.c {
		database := g.DB[dbConfig.Driver]
		if database == nil {
			panic(fmt.Sprintf("The database type %s is currently not supported. The database name is %s", dbConfig.Driver, dbConfig.Database))
		}
		conn, err := database.Conn(dbConfig)
		if err != nil {
			panic(fmt.Sprintf("The database %s connection failed， error: %s", dbConfig.Database, err))
		}

		if dbConfig.Alias != "" {
			databases[dbConfig.Alias] = conn
		} else {
			databases[dbConfig.Database] = conn
		}

		sqlDB, err := conn.DB()

		// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)

		if dbConfig.Alias == "default" {
			databases[dbConfig.Database] = conn
		}
	}

	g.client = databases
}

func (g *Gorm) ConnType(database string) bool {
	if g.DB[database] == nil {
		return false
	}
	return true
}

func (g *Gorm) Result(c *databese.Engine) {
	c.Gorm = g.client
}
