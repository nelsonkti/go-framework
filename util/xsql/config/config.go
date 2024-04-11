package config

import (
	"fmt"
	"go-framework/util/helper"
)

// xsql config 配置信息
type DBConfig struct {
	Driver   string   `json:"driver"`
	Host     string   `json:"host"`
	Sources  []string `json:"sources"`
	Replicas []string `json:"replicas"`
	Port     int      `json:"port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Database string   `json:"database"`
	Alias    string   `json:"alias"`
}

func Marshal(v interface{}) []byte {
	marshal, err := helper.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("connection database: %s", err))
	}

	return marshal
}
