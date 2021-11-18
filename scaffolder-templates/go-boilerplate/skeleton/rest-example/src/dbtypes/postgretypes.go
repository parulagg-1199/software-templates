package dbtypes

import (
	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"

	"git.xenonstack.com/lib/golang-boilerplate/db"

	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/methods"
)

// User is structure which is used to store user information
type User struct {
	Email    string `json:"email" gorm:"unique_index;"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

// CreatePostgreDBTablesIfNotExists is method which is used to create or update table in postgres database
func CreatePostgreDBTablesIfNotExists() {
	// create postgre database if database is not initialized
	sql := db.NewPostgresql()
	sql = config.Conf.Postgres
	err := sql.CreatePostgresDatabase()
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	zap.S().Info("postgres database had been created")
	// connecting db using connection string
	str, err := config.PostgreString()
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	// initialize db connection
	dbClient, err := gorm.Open("postgres", str)
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	// close db instance whenever whole work completed
	defer dbClient.Close()

	//create tables
	if !(dbClient.HasTable(User{})) {
		dbClient.CreateTable(User{})

		//creating admin account
		adminAcc := initAdminAccount()
		dbClient.Create(&adminAcc)
	}

	dbClient.AutoMigrate(&User{})

	zap.S().Info("postgres database initialized successfully")
}

func initAdminAccount() User {
	// fetching info from env variables
	adminEmail := config.Conf.Admin.Email
	if adminEmail == "" {
		adminEmail = "admin@xenonstack.com"
	}
	adminPass := config.Conf.Admin.Pass
	if adminPass == "" {
		adminPass = "admin"
	}
	// return struct with details of admin
	return User{Password: methods.HashForNewPassword(adminPass),
		Email: adminEmail,
		Role:  "admin"}
}
