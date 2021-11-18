package api

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/info"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SaveInformation is an api handler to save information of user in mongodb
func SaveInformation(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "save user information")

	span.LogKV("task", "binding body data")
	//bind details from body
	var data dbtypes.UserInfo
	if err := c.BindJSON(&data); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": "Please pass name and contact no."})
		return
	}

	span.LogKV("task", "extract jwt claims")
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	span.LogKV("task", "extract email data from jwt claims")
	// extract email from jwt claims and before assigning check email exists in claims map
	val, ok := claims["email"]
	if ok {
		data.Email = val.(string)
	} else {
		span.LogKV("task", "send final output")
		zap.S().Error("email claim is not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again after some time"})
		return
	}
	span.LogKV("task", "save data in db")
	err := info.Save(data)
	span.LogKV("task", "send final output")
	if err != nil {
		zap.S().Error(err.Error())
		c.JSON(500, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"error": false, "message": "Information saved succesfully"})
	return
}

// FetchInformation is an api handler for fetching information related to user from mongodb
func FetchInformation(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "fetch user information")

	span.LogKV("task", "extract jwt claims")
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	span.LogKV("task", "extract email data from jwt claims")
	// extract email from jwt claims and before assigning check email exists in claims map
	val, ok := claims["email"]
	if ok {
		span.LogKV("task", "fetch details")
		data, err := info.Fetch(val.(string))
		span.LogKV("task", "send final output")
		if err != nil {
			c.JSON(500, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"error": false, "information": data})
	} else {
		span.LogKV("task", "send final output")
		zap.S().Error("email claim is not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again after some time"})
		return
	}
}
