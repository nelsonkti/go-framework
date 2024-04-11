package config

type Conf struct {
	App    App           `json:"app"`
	Server Server        `json:"server"`
	DB     map[string]DB `json:"db"`
	Redis  []Redis       `json:"redis"`
	MQ     MQ            `json:"mq"`
}

type App struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

type Server struct {
	Http Network `json:"http"`
	Rpc  Network `json:"rpc"`
}

type Network struct {
	Addr string `json:"addr"`
}

type DB struct {
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

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database int    `json:"database"`
	Alias    string `json:"alias"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type MQ struct {
	Endpoint  []string `json:"endpoint"`
	AccessKey string   `json:"access_key"`
	SecretKey string   `json:"secret_key"`
	Namespace string   `json:"namespace"`
	Env       string   `json:"env"`
}
