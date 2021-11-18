package api

import (
	"fmt"

	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Data is structure for binding data from request body of request to save data in redis
type Data struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// SaveData is an api handler for saving data in redis database
func SaveData(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "save key value in redis")

	//bind details from body
	span.LogKV("task", "binding body data")
	var data Data
	if err := c.BindJSON(&data); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": "Please pass key and value"})
		return
	}

	// save data in redis by passing key and value
	span.LogKV("task", "save data in redis db")
	err := config.Conf.Redis.SaveData(data.Key, data.Value)
	span.LogKV("task", "send final output")
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "data successfully saved in redis db",
	})
}

// CheckKey is an api hadnler for checking key exists in redis database
func CheckKey(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "fetch value of key from redis")

	// fetch key from request query data
	span.LogKV("task", "fetching key from query data")
	key := c.Query("key")
	zap.S().Info(key)
	if key == "" {
		span.LogKV("task", "send final output")
		c.JSON(400, gin.H{
			"error":   true,
			"message": "please pass key in query",
		})
		return
	}
	// fetch value correspondance to key from redis
	span.LogKV("task", "check key exists")
	val, err := config.Conf.Redis.CheckData(key)
	span.LogKV("task", "send final output")
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "no key exists",
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "key is there and its value is " + fmt.Sprint(val),
	})
}

// DeleteKey is an api handler for deleting data from redis on basis of key
func DeleteKey(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "delete key from redis")

	// fetch key from request query data
	span.LogKV("task", "fetching key from query data")
	key := c.Query("key")
	zap.S().Info(key)
	if key == "" {
		span.LogKV("task", "send final output")
		c.JSON(400, gin.H{
			"error":   true,
			"message": "please pass key in query",
		})
		return
	}
	// delete data from redis by passing key
	span.LogKV("task", "delete key")
	err := config.Conf.Redis.DeleteData(key)
	span.LogKV("task", "send final output")
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "no key exists",
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "key successfully deleted from redis",
	})
}
