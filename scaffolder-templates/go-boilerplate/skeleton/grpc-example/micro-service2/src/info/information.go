package info

import (
	"errors"

	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/config"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service2/src/types"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

// Save is a function to save user information in database
func Save(parentspan opentracing.Span, data types.Info) error {
	// start span from parent context
	span := opentracing.StartSpan("save user information in mongodb", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "start mongodb session and select collection")
	//create mongo session
	session, err := config.MongoSession()
	if err != nil {
		span.LogKV("task", "return due to mongodb error")
		zap.S().Error("error connecting DB", err)
		return err
	}
	defer session.Close()
	//select collection
	collection := session.DB(config.Conf.Mongo.DBName).C("user")

	//insert data in collection
	span.LogKV("task", "insert data in collection")
	err = collection.Insert(&data)
	span.LogKV("task", "return error if any")
	return err
}

// Fetch is a function to get user information from database
func Fetch(parentspan opentracing.Span, email string) (types.Info, error) {
	// start span from parent context
	span := opentracing.StartSpan("save user information in mongodb", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "start mongodb session and select collection")
	//create mongo session
	session, err := config.MongoSession()
	if err != nil {
		span.LogKV("task", "return due to mongodb error")
		zap.S().Error("error connecting DB", err)
		return types.Info{}, err
	}
	defer session.Close()
	//select collection
	collection := session.DB(config.Conf.Mongo.DBName).C("user")

	//fetch data from database
	span.LogKV("task", "fetch data from collection")
	data := []types.Info{}
	err = collection.Find(bson.M{"email": email}).All(&data)
	span.LogKV("task", "return value")
	if err != nil {
		return types.Info{}, err
	}
	if len(data) == 0 {
		return types.Info{}, errors.New("there is no information saved")
	}
	return data[0], nil
}
