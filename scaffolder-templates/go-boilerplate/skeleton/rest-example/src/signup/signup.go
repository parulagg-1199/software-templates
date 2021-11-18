package signup

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"

	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"errors"
)

// Signup is a method used for gorm there is need to add a blank import for dialects
func Signup(user dbtypes.User) (int, error) {
	str, err := config.PostgreString()
	if err != nil {
		zap.S().Error(err)
		return 500, err
	}
	// connecting to db
	db, err := gorm.Open("postgres", str)
	if err != nil {
		zap.S().Error(err)
		return 500, err
	}
	// close db instance whenever whole work completed
	defer db.Close()

	//save information in database
	row := db.Create(&user).RowsAffected
	zap.S().Info(row)
	if row == 0 {
		return 400, errors.New("Email already exists")
	}
	return 200, nil
}

// DeleteUser is a method to delete user information from portgre database
func DeleteUser(parentspan opentracing.Span, email string) error {
	span := opentracing.StartSpan("delete account from postgre", opentracing.ChildOf(parentspan.Context()))
	zap.S().Info(span.BaggageItem("email"))
	defer span.Finish()
	span.LogKV("task", "intialise db connection")
	str, err := config.PostgreString()
	if err != nil {
		span.LogKV("task", "send final output")
		zap.S().Error(err)
		return err
	}
	// connecting to db
	db, err := gorm.Open("postgres", str)
	if err != nil {
		span.LogKV("task", "send final output")
		zap.S().Error(err)
		return err
	}
	// close db instance whenever whole work completed
	defer db.Close()

	span.LogKV("task", "delete user account from db")
	row := db.Where("email=?", email).Delete(dbtypes.User{}).RowsAffected
	span.LogKV("task", "send final output")
	if row == 0 {
		return errors.New("no user exists")
	}
	return nil
}
