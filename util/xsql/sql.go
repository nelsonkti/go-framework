package xsql

import (
	"encoding/json"
	"go-framework/util/xsql/config"
	"go-framework/util/xsql/databese"
)

func NewClient(c interface{}) (*databese.Engine, error) {
	cByte, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	databases := make(map[string]config.DBConfig)
	err = json.Unmarshal(cByte, &databases)
	if err != nil {
		return nil, err
	}

	engine, err := db.NewDB(databases).InitDatabases()
	if err != nil {
		return nil, err
	}

	return engine, nil
}
