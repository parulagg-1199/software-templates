package config

import (
	"log"
	"os"
	"time"

	"git.xenonstack.com/lib/golang-boilerplate/db"
	"git.xenonstack.com/lib/golang-boilerplate/redis"

	"github.com/BurntSushi/toml"
)

// Config is a structure for configuration
type Config struct {
	Postgres db.PostgreConfig
	Admin    Admin
	Redis    redis.Config
	Service  Service
	JWT      JWT
	Jaeger   Tracer
}

// Redis is a structure for redis database configuration
type Redis struct {
	Database   string
	Host       string
	Port       string
	Pass       string
	ExpireTime time.Duration
}

// Service is a structure for service specific related configuration
type Service struct {
	Port        string
	Environment string
}

// JWT is structure for jwt token specific configuration
type JWT struct {
	PrivateKey string
	ExpireTime time.Duration
}

// Tracer is a strcuture for jaeger configuration
type Tracer struct {
	Host string
	Port string
}

// Admin is a structure for admin account credentials
type Admin struct {
	Email string
	Pass  string
}

// Conf is a global variable for configuration
var Conf Config

// TomlFile is a global variable for toml file path
var TomlFile string

// ConfigurationWithEnv is a method to initialize configuration with environment variables
func ConfigurationWithEnv() {
	// Postgres database configuration
	Conf.Postgres.Host = os.Getenv("POSTGRESQL_DB_HOST")
	Conf.Postgres.Port = os.Getenv("POSTGRESQL_DB_PORT")
	Conf.Postgres.User = os.Getenv("POSTGRESQL_DB_USER")
	Conf.Postgres.Pass = os.Getenv("POSTGRESQL_DB_PASS")
	Conf.Postgres.DBName = os.Getenv("POSTGRESQL_DB_NAME")

	// redis database configuration
	Conf.Redis.Database = os.Getenv("REDIS_DB")
	Conf.Redis.Host = os.Getenv("REDIS_HOST")
	Conf.Redis.Port = os.Getenv("REDIS_PORT")
	Conf.Redis.Pass = os.Getenv("REDIS_PASS")

	// admin account credentials configuration
	Conf.Admin.Email = os.Getenv("GRAPHQL_EXAMPLE_ADMIN_EMAIL")
	Conf.Admin.Pass = os.Getenv("GRAPHQL_EXAMPLE_ADMIN_PASS")

	// jaeger specific configuration
	Conf.Jaeger.Host = os.Getenv("JAEGER_AGENT_HOST")
	Conf.Jaeger.Port = os.Getenv("JAEGER_AGENT_PORT")

	// if service port is not defined set default port
	if os.Getenv("GRAPHQL_EXAMPLE_HTTP_PORT") != "" {
		Conf.Service.Port = os.Getenv("GRAPHQL_EXAMPLE_HTTP_PORT")
	} else {
		Conf.Service.Port = "8000"
	}
	Conf.Service.Environment = os.Getenv("ENVIRONMENT")

	//service specific configuration
	Conf.JWT.PrivateKey = os.Getenv("GRAPHQL_EXAMPLE_PRIVATE_KEY")

	// constants
	//JWT Token Timeout in minutes
	Conf.JWT.ExpireTime = time.Minute * 30
	Conf.Redis.ExpireTime = time.Minute * 5

}

// ConfigurationWithToml is a method to initialize configuration with toml file
func ConfigurationWithToml(filePath string) {
	// set varible as file path if configuration is done using toml
	TomlFile = filePath

	// parse toml file and save data config structure
	_, err := toml.DecodeFile(filePath, &Conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// if service port is not defined set default port
	if Conf.Service.Port == "" {
		Conf.Service.Port = "8000"
	}
	// constants
	//JWT Token Timeout in minutes
	Conf.JWT.ExpireTime = time.Minute * 30
	Conf.Redis.ExpireTime = time.Minute * 5
}

// SetConfig is a method to re-intialise configuration at runtime
func SetConfig() {
	if TomlFile == "" {
		ConfigurationWithEnv()
	} else {
		ConfigurationWithToml(TomlFile)
	}
}

// PostgreString is a method which return PostgreSql connection string
func PostgreString() (string, error) {
	SetConfig()
	return Conf.Postgres.PostgreString()
}
