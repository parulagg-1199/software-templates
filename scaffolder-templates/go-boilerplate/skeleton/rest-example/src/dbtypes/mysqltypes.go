package dbtypes

import (
	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"

	"git.xenonstack.com/lib/golang-boilerplate/db"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"
)

// UserDetail is structure which is used to store user details
type UserDetail struct {
	Email   string `json:"-" gorm:"unique_index;"`
	Name    string `json:"name" binding:"required"`
	Contact string `json:"contact" binding:"required"`
}

// CreateMySQLDBTablesIfNotExists is method which is used to create or update table in mysql database
func CreateMySQLDBTablesIfNotExists() {
	// create postgre database if database is not initialized
	sql := db.NewMysql()
	sql = config.Conf.Mysql
	err := sql.CreateMysqlDatabase()
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	zap.S().Info("mysql database had been created")
	// connecting db using connection string
	str, err := config.MysqlString()
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	// initialize db connection
	dbClient, err := gorm.Open("mysql", str)
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	// close db instance whenever whole work completed
	defer dbClient.Close()

	//create tables
	if !(dbClient.HasTable(UserDetail{})) {
		dbClient.CreateTable(UserDetail{})
	}

	dbClient.AutoMigrate(&UserDetail{})

	zap.S().Info("mysql database initialized successfully")
}
