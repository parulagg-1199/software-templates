package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MysqlConfig is a structure for configuration of mysql
type MysqlConfig struct {
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
}

// NewMysql is a method which returns deafult config of mysql
func NewMysql() MysqlConfig {
	return MysqlConfig{
		Host:   "localhost",
		Port:   "3306",
		User:   "root",
		Pass:   "",
		DBName: "default",
	}
}

// MysqlString is a method which return mysql connection string
func (config *MysqlConfig) MysqlString() (string, error) {
	if config.Host == "" {
		return "", errors.New("please set mysql host")
	}
	if config.Port == "" {
		return "", errors.New("please set mysql port")
	}
	if config.User == "" {
		return "", errors.New("please set mysql user")
	}
	if config.DBName == "" {
		return "", errors.New("please set mysql dbname")
	}

	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.DBName)

	return str, nil
}

// CreateMysqlDatabase is a method which is used for creating mysql database
func (config *MysqlConfig) CreateMysqlDatabase() error {
	// connecting with mysql root db
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		config.User,
		config.Pass,
		config.Host,
		config.Port))
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	// executing create database query.
	db.Exec("create database " + config.DBName + ";")
	return nil
}

// MysqlHealthz is a method to check health of mysql
func (config *MysqlConfig) MysqlHealthz() bool {
	//checking health of mysql
	// connecting to db
	str, err := config.MysqlString()
	if err != nil {
		log.Println(err)
		return false
	}
	db, err := gorm.Open("mysql", str)
	if err != nil {
		log.Println(err)
		return false
	}
	// close db instance whenever whole work completed
	defer db.Close()

	return true
}
