package redis

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

// Config is a structure for redis configuration
type Config struct {
	Database   string
	Host       string
	Port       string
	Pass       string
	ExpireTime time.Duration
}

// New is method which return default setting of redis configuration
func New() Config {
	return Config{
		Database:   "1",
		Host:       "localhost",
		Port:       "6379",
		Pass:       "",
		ExpireTime: time.Minute * 1,
	}
}

func (config *Config) initialise() error {
	db, _ := strconv.Atoi(config.Database)
	// redis db client creations
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: fmt.Sprintf("%s", config.Pass),
		DB:       db,
	})

	// check connection with server
	pong, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(pong)
	return nil
}

// SaveData is method to save and update data in redis database
func (config *Config) SaveData(key string, value interface{}) error {
	err := config.initialise()
	if err != nil {
		log.Println(err)
		return err
	}
	err = client.Set(key, value, config.ExpireTime).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CheckData is a method check key is there in redis database
func (config *Config) CheckData(key string) (interface{}, error) {
	err := config.initialise()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	val, err := client.Get(key).Result()
	if err != nil {
		// when token not exist
		log.Println(err)
		return nil, err
	}
	log.Println(val)
	return val, nil
}

// DeleteData is a method delete data from database
func (config *Config) DeleteData(key string) error {
	err := config.initialise()
	if err != nil {
		log.Println(err)
		return err
	}
	val, err := client.Del(key).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	if val == 0 {
		return errors.New("No key exists")
	}
	return nil
}
