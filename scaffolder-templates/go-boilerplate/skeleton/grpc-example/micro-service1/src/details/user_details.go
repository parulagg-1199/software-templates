package details

import (
	"strconv"

	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/config"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/src/dbtypes"

	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	// for gorm there is need to add a blank import for dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func checkData(parentspan opentracing.Span, data dbtypes.UserDetail) error {
	span := opentracing.StartSpan("validate user detail data", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "check name")
	//check name not contains any number or special characters
	if strings.ContainsAny(data.Name, "1234567890{}()~`:;<>,./!@#$%^_+=[]|&*-?\"\\'") {
		return errors.New("Name can contain only alphabet and space")
	}
	span.LogKV("task", "return output")
	return nil
}

// Save is a method to save user details in database
func Save(parentspan opentracing.Span, data dbtypes.UserDetail) (string, error) {
	// start span from parent context
	span := opentracing.StartSpan("save user details in mysql db", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "check data validation")
	//check data validation
	err := checkData(span, data)
	if err != nil {
		zap.S().Error(err)
		span.LogKV("task", "return output")
		return err.Error(), err
	}

	span.LogKV("task", "intialise db connection")
	str, err := config.MysqlString()
	if err != nil {
		zap.S().Error(err)
		span.LogKV("task", "return output")
		return err.Error(), err
	}
	// connecting to db
	db, err := gorm.Open("mysql", str)
	if err != nil {
		zap.S().Error(err)
		span.LogKV("task", "return output")
		return "Unable to connect to database", err
	}
	// close db instance whenever whole work completed
	defer db.Close()

	span.LogKV("task", "fetch user details from db")
	// fetch user information
	details := []dbtypes.UserDetail{}
	db.Where("email=?", data.Email).Find(&details)
	if len(details) != 0 {
		span.LogKV("task", "update user details in db if user data is already there")
		// if user already exists
		row := db.Exec("update user_details set name='" + data.Name + "', contact=" + strconv.Itoa(int(data.Contact)) + " where email='" + data.Email + "';").RowsAffected
		zap.S().Info(row)
		span.LogKV("task", "return output")
		return data.Name + " your profile get updated", nil
	}
	span.LogKV("task", "save user details in db if user data is not already there")
	// if no details had been there for user
	db.Create(&data)
	span.LogKV("task", "return output")
	return data.Name + " your personal details had been saved", nil

}

// Fetch is a method to get user details from database
func Fetch(parentspan opentracing.Span) (dbtypes.UserDetail, error) {
	// start span from parent context
	span := opentracing.StartSpan("fetch user details from mysql", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "intialise db connection")
	str, err := config.MysqlString()
	if err != nil {
		zap.S().Error(err)
		return dbtypes.UserDetail{}, err
	}
	// connecting to db
	db, err := gorm.Open("mysql", str)
	if err != nil {
		zap.S().Error(err)
		return dbtypes.UserDetail{}, err
	}
	// close db instance whenever whole work completed
	defer db.Close()

	//fetch email from baggage
	email := span.BaggageItem("email")
	zap.S().Info("email in baggage===", email)
	span.LogKV("task", "fetch user details from db")
	// fetch user information
	details := []dbtypes.UserDetail{}
	db.Where("email=?", email).Find(&details)
	zap.S().Info(len(details))
	span.LogKV("task", "send final output")
	if len(details) == 0 {
		return dbtypes.UserDetail{}, errors.New("No details found")
	}
	return details[0], nil
}

// Delete is a method to delete user details from database
func Delete(parentspan opentracing.Span, email string) error {
	// start span from parent context
	span := opentracing.StartSpan("delete user details from mysql", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "intialise db connection")
	str, err := config.MysqlString()
	if err != nil {
		span.LogKV("task", "send final output")
		zap.S().Error(err)
		return err
	}
	// connecting to db
	db, err := gorm.Open("mysql", str)
	if err != nil {
		span.LogKV("task", "send final output")
		zap.S().Error(err)
		return err
	}
	// close db instance whenever whole work completed
	defer db.Close()

	span.LogKV("task", "delete user account from db")
	// delete data form database
	row := db.Where("email=?", email).Delete(dbtypes.UserDetail{}).RowsAffected
	span.LogKV("task", "send final output")
	if row == 0 {
		return errors.New("no user details exists")
	}
	return nil
}
