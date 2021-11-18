package db

import (
	mgo "gopkg.in/mgo.v2"

	"log"
	"time"
)

// MongoConfig is a structure for configuration of mongo database
type MongoConfig struct {
	Host   string
	Auth   string
	User   string
	Pass   string
	DBName string
}

// NewMongo is a method which returns deafult config of mongodb
func NewMongo() MongoConfig {
	return MongoConfig{
		Host:   "localhost:27017",
		Auth:   "admin",
		User:   "",
		Pass:   "",
		DBName: "default",
	}
}

// MongoDBConfig is a method which returns mongo dial info structure  object
func (config *MongoConfig) MongoDBConfig() mgo.DialInfo {
	return mgo.DialInfo{
		Addrs:    []string{config.Host},
		Timeout:  60 * time.Second,
		Database: config.Auth,
		Username: config.User,
		Password: config.Pass,
	}
}

// MongoHealthz is a method to check health of mongo database
func (config *MongoConfig) MongoHealthz() bool {
	// checking health of mongodb
	conf := config.MongoDBConfig()
	//creating session using dial info
	session, err := mgo.DialWithInfo(&conf)
	if err != nil {
		log.Println("error connecting DB", err)
		return false
	}
	//close session
	defer session.Close()
	return true
}
