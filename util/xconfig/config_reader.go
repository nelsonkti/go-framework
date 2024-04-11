package xconfig

type ConfigReader interface {
	Load() (map[string]interface{}, error)
}
