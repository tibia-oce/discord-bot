package configs

import (
	"fmt"
)

const (
	EnvDBHostKey = "MYSQL_HOST"
	EnvDBNameKey = "MYSQL_DBNAME"
	EnvDBPassKey = "MYSQL_PASS"
	EnvDBPortKey = "MYSQL_PORT"
	EnvDBUserKey = "MYSQL_USER"
)

type DBConfigs struct {
	Host string
	Port int
	Name string
	User string
	Pass string
	Config
}

func (dbConfigs *DBConfigs) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		dbConfigs.User,
		dbConfigs.Pass,
		dbConfigs.Host,
		dbConfigs.Port,
		dbConfigs.Name,
	)
}

func (dbConfigs *DBConfigs) format() string {
	return fmt.Sprintf(
		"Database: %s:%d/%s",
		dbConfigs.Host,
		dbConfigs.Port,
		dbConfigs.Name,
	)
}
func GetDBConfigs() DBConfigs {
	return DBConfigs{
		Host: GetEnvStr(EnvDBHostKey, "database"),
		Port: GetEnvInt(EnvDBPortKey, 3306),
		Name: GetEnvStr(EnvDBNameKey, "forgottenserver"),
		User: GetEnvStr(EnvDBUserKey, "forgottenserver"),
		Pass: GetEnvStr(EnvDBPassKey, "forgottenserver"),
	}
}
