package postgredb

import (
	"sync"

	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	"moul.io/zapgorm"

	"git.xenonstack.com/lib/golang-boilerplate/db"
	"git.xenonstack.com/lib/golang-boilerplate/logger"

	"git.xenonstack.com/util/golang-boilerplate/graphql-example/config"
	"git.xenonstack.com/util/golang-boilerplate/graphql-example/src/methods"
	"git.xenonstack.com/util/golang-boilerplate/graphql-example/src/models"
)

var once sync.Once

// Connect is function for initializing connection with postgres database
func Connect() (*gorm.DB, error) {
	var dbClient *gorm.DB
	// db connection string
	str, err := config.PostgreString()
	if err != nil {
		zap.S().Error(err)
		return dbClient, err
	}
	once.Do(func() {
		// initialise dbclient
		dbClient, err = gorm.Open("postgres", str)
		// enable debug mode
		dbClient = dbClient.Debug()
		// change logger to zap
		dbClient.SetLogger(zapgorm.New(logger.Log))
		zap.S().Info(err)
		return
	})
	return dbClient, err
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
	db, err := gorm.Open("postgres", str)
	if err != nil {
		zap.S().Error(err.Error())
		return
	}
	// close db instance whenever whole work completed
	defer db.Close()

	//create tables
	if !(db.HasTable(models.UserDetail{})) {
		db.CreateTable(models.UserDetail{})

		// creating admin account
		adminAcc := initAdminAccount()
		db.Create(&adminAcc)
	}

	if !(db.HasTable(models.Address{})) {
		db.CreateTable(models.Address{})
	}

	db.AutoMigrate(&models.UserDetail{}, &models.Address{})

	zap.S().Info("postgres database initialized successfully")
}

func initAdminAccount() models.UserDetail {
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
	return models.UserDetail{
		Password: methods.HashForNewPassword(adminPass),
		Email:    adminEmail,
		Role:     "admin",
	}
}
