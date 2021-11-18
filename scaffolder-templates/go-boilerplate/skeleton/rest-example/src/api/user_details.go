package api

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/details"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddUserDetails is an api handler to save user details in mysql
func AddUserDetails(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "save user details")

	span.LogKV("task", "binding body data")
	//bind details from body
	var data dbtypes.UserDetail
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
	code, err, msg := details.Save(span, data)
	span.LogKV("task", "send final output")
	if err != nil {
		zap.S().Error(err.Error())
		c.JSON(code, gin.H{"error": true, "message": msg})
		return
	}
	c.JSON(code, gin.H{"error": false, "message": msg})
	return
}

// GetUserDetails is an api handler to fetch user details from mysql
func GetUserDetails(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "fetch user details")

	span.LogKV("task", "extract jwt claims")
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	span.LogKV("task", "extract email data from jwt claims")
	// extract email from jwt claims and before assigning check email exists in claims map
	val, ok := claims["email"]
	if ok {
		span.SetBaggageItem("email", val.(string))
		span.LogKV("task", "fetch details")
		data, err := details.Fetch(span)
		span.LogKV("task", "send final output")
		if err != nil {
			c.JSON(500, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"error": false, "details": data})
	} else {
		span.LogKV("task", "send final output")
		zap.S().Error("email claim is not set")
		c.JSON(500, gin.H{"error": true, "message": "Please login again after some time"})
		return
	}
}
