package config

type MysqlConfig struct {
	Host     string
	Port     int
	Schema   string
	Username string
	Password string
}

func NewMysqlConfig(host string, port int, schema string, username string, password string) *MysqlConfig {
	return &MysqlConfig{host, port, schema, username, password}
}
