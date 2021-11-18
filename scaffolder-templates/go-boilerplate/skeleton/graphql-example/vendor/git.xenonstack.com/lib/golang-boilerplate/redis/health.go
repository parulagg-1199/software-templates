package redis

import (
	"log"
)

// Healthz is a method to check connection with redis database
func (config *Config) Healthz() bool {
	// initialise client
	err := config.initialise()
	if err != nil {
		log.Println(err)
		return false
	}
	// save key-value pair in redis
	err = client.Set("key", "value", config.ExpireTime).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	// get above key from redis
	val, err := client.Get("key").Result()
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("key", val)
	return val == "value"
}
