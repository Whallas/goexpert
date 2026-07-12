package configs

import (
	"fmt"
	"os"
)

type Conf struct {
	DBDriver          string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	WebServerPort     string
	GRPCServerPort    string
	GraphQLServerPort string
}

func LoadConfig() *Conf {
	return &Conf{
		DBDriver:          getEnv("DB_DRIVER", "mysql"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", "root"),
		DBName:            getEnv("DB_NAME", "orders"),
		WebServerPort:     getEnv("WEB_SERVER_PORT", "8000"),
		GRPCServerPort:    getEnv("GRPC_SERVER_PORT", "50051"),
		GraphQLServerPort: getEnv("GRAPHQL_SERVER_PORT", "8080"),
	}
}

func (c *Conf) DBSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Conf) MigrateURL() string {
	return fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
