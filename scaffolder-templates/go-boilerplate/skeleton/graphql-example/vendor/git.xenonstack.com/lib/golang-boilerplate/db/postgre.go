package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PostgreConfig is a structure for configuration of postgres
type PostgreConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

// NewPostgresql is a method which returns deafult config of postgre
func NewPostgresql() PostgreConfig {
	return PostgreConfig{
		Host:   "localhost",
		Port:   "5432",
		User:   "postgres",
		Pass:   "1234",
		DBName: "default",
	}
}

// PostgreString is a method which return postgres connection string
func (config *PostgreConfig) PostgreString() (string, error) {
	if config.Host == "" {
		return "", errors.New("please set postgres host")
	}
	if config.Port == "" {
		return "", errors.New("please set postgres port")
	}
	if config.User == "" {
		return "", errors.New("please set postgres user")
	}
	if config.DBName == "" {
		return "", errors.New("please set postgres dbname")
	}

	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Pass,
		config.DBName)

	return str, nil
}

// CreatePostgresDatabase is method is used for creating postgres database
func (config *PostgreConfig) CreatePostgresDatabase() error {
	// connecting with postgres root db
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Pass,
		"postgres", "disable"))
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	// executing create database query.
	db.Exec("create database " + config.DBName + ";")
	return nil
}

// PostgresHealthz is a method to check health of postgres
func (config *PostgreConfig) PostgresHealthz() bool {
	//checking health of postgres sql
	str, err := config.PostgreString()
	if err != nil {
		log.Println(err)
		return false
	}
	// connecting to db
	db, err := gorm.Open("postgres", str)
	if err != nil {
		log.Println(err)
		return false
	}
	// close db instance whenever whole work completed
	defer db.Close()

	return true
}
