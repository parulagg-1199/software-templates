package info

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"github.com/opentracing/opentracing-go"

	"errors"

	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Save is a method to store information related to user in mongodb
func Save(data dbtypes.UserInfo) error {
	//create mongo session
	session, err := config.MongoSession()
	if err != nil {
		zap.S().Error("error connecting DB", err)
		return err
	}
	defer session.Close()
	//select collection
	collection := session.DB(config.Conf.Mongo.DBName).C("user")

	//insert data in collection
	err = collection.Insert(&data)
	if err == nil {
		//create index
		ensureIndex()
	}
	return err
}

// Fetch is a method to fetch information related to user from mongodb
func Fetch(email string) (dbtypes.UserInfo, error) {
	//create mongo session
	session, err := config.MongoSession()
	if err != nil {
		zap.S().Error("error connecting DB", err)
		return dbtypes.UserInfo{}, err
	}
	defer session.Close()
	//select collection
	collection := session.DB(config.Conf.Mongo.DBName).C("user")

	data := []dbtypes.UserInfo{}
	//insert data in collection
	err = collection.Find(bson.M{"Email": email}).All(&data)
	if err != nil {
		zap.S().Error(err)
		return dbtypes.UserInfo{}, err
	}
	if len(data) == 0 {
		return dbtypes.UserInfo{}, errors.New("No information found")
	}
	return data[0], nil
}

// Delete is a method to delete user information related to user from mongodb
func Delete(parentspan opentracing.Span, email string) error {
	span := opentracing.StartSpan("delete user information from mongodb", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("task", "intialise db connection and select collection")
	//create mongo session
	session, err := config.MongoSession()
	if err != nil {
		zap.S().Error("error connecting DB", err)
		return err
	}
	defer session.Close()
	//select collection
	collection := session.DB(config.Conf.Mongo.DBName).C("user")

	span.LogKV("task", "delete user information from db")
	err = collection.Remove(bson.M{"Email": email})
	span.LogKV("task", "send final output")
	return err
}

// ensureIndex is a method to create index in mongodb
func ensureIndex() {
	//fetch mongodb config
	//create mongo session
	session, err := config.MongoSession()
	if err != nil {
		zap.S().Error("error connecting DB", err)
		return
	}
	defer session.Close()
	//select collection
	collection := session.DB(config.Conf.Mongo.DBName).C("user")
	//create collection on email field
	index := mgo.Index{
		Key:      []string{"Email"},
		Unique:   true,
		DropDups: true,
	}
	//delete previous index
	err = collection.DropIndex("Email")
	if err != nil {
		zap.S().Error(err)
	}
	//create new index
	err = collection.EnsureIndex(index)
	if err != nil {
		zap.S().Error(err)
	}
}
