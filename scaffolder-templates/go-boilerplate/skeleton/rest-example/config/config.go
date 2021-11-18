package config

import (
	"log"
	"os"
	"time"

	"git.xenonstack.com/lib/golang-boilerplate/db"
	"git.xenonstack.com/lib/golang-boilerplate/redis"

	"github.com/BurntSushi/toml"
	mgo "gopkg.in/mgo.v2"
)

// Config is a structure for configuration
type Config struct {
	Postgres db.PostgreConfig
	Mysql    db.MysqlConfig
	Mongo    db.MongoConfig
	Redis    redis.Config
	Service  Service
	Admin    Admin
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

// Admin is a structure for admin account credentials
type Admin struct {
	Email string
	Pass  string
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

	// mysql database configuration
	Conf.Mysql.Host = os.Getenv("MYSQL_DB_HOST")
	Conf.Mysql.Port = os.Getenv("MYSQL_DB_PORT")
	Conf.Mysql.User = os.Getenv("MYSQL_DB_USER")
	Conf.Mysql.Pass = os.Getenv("MYSQL_DB_PASS")
	Conf.Mysql.DBName = os.Getenv("MYSQL_DB_NAME")

	// cockroach database configuration
	Conf.Mongo.Host = os.Getenv("MONGO_DB_HOST")
	Conf.Mongo.Auth = os.Getenv("MONGO_AUTH_DB_NAME")
	Conf.Mongo.User = os.Getenv("MONGO_DB_USER")
	Conf.Mongo.Pass = os.Getenv("MONGO_DB_PASS")
	Conf.Mongo.DBName = os.Getenv("MONGO_DB_NAME")

	// redis database configuration
	Conf.Redis.Database = os.Getenv("REDIS_DB")
	Conf.Redis.Host = os.Getenv("REDIS_HOST")
	Conf.Redis.Port = os.Getenv("REDIS_PORT")
	Conf.Redis.Pass = os.Getenv("REDIS_PASS")

	// admin account credentials configuration
	Conf.Admin.Email = os.Getenv("REST_EXAMPLE_ADMIN_EMAIL")
	Conf.Admin.Pass = os.Getenv("REST_EXAMPLE_ADMIN_PASS")

	// jaeger specific configuration
	Conf.Jaeger.Host = os.Getenv("JAEGER_AGENT_HOST")
	Conf.Jaeger.Port = os.Getenv("JAEGER_AGENT_PORT")

	// if service port is not defined set default port
	if os.Getenv("REST_EXAMPLE_HTTP_PORT") != "" {
		Conf.Service.Port = os.Getenv("REST_EXAMPLE_HTTP_PORT")
	} else {
		Conf.Service.Port = "8000"
	}
	Conf.Service.Environment = os.Getenv("ENVIRONMENT")

	//service specific configuration
	Conf.JWT.PrivateKey = os.Getenv("REST_EXAMPLE_PRIVATE_KEY")

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

// MysqlString is a method which return mysql connection string
func MysqlString() (string, error) {
	SetConfig()
	return Conf.Mysql.MysqlString()
}

// PostgreString is a method which return PostgreSql connection string
func PostgreString() (string, error) {
	SetConfig()
	return Conf.Postgres.PostgreString()
}

// MongoDialInfo is a function which return mongo dial info when need authorization
func MongoDialInfo() mgo.DialInfo {
	return Conf.Mongo.MongoDBConfig()
}

// MongoSession is a method which create session with mongo database
func MongoSession() (*mgo.Session, error) {
	SetConfig()
	var session *mgo.Session
	var err error
	if Conf.Mongo.User == "" && Conf.Mongo.Pass == "" {
		// session without authorization
		session, err = mgo.Dial(Conf.Mongo.Host)
		if err != nil {
			log.Println(err)
			return session, err
		}
	} else {
		config := MongoDialInfo()
		// session with authorization
		session, err = mgo.DialWithInfo(&config)
		if err != nil {
			log.Println(err)
			return session, err
		}
	}
	return session, nil
}
