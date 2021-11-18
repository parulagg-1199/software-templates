package api

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/details"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/info"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/signup"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeleteAccount is an api handler to delete all data related to user on basis of email from databases
func DeleteAccount(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "delete a user by admin only")

	// fetch email from query data from request
	span.LogKV("task", "fetching email from query data")
	email := c.Query("email")
	zap.S().Info(email)
	if email == "" {
		span.LogKV("task", "send final output")
		c.JSON(400, gin.H{
			"error":   true,
			"message": "please pass email in query",
		})
		return
	}
	span.SetBaggageItem("email", email)
	// delete user infromation from postgres
	span.LogKV("task", "delete user account from postgre using email Id")
	err := signup.DeleteUser(span, email)
	zap.S().Info("account error ", err)
	if err != nil {
		span.LogKV("task", "send final output")
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	// delete user details from mysql
	span.LogKV("task", "delete user details from mysql using email Id")
	err = details.Delete(span, email)
	zap.S().Info("details error ", err)
	if err != nil {
		span.LogKV("task", "send final output")
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	// delete infromation related to user from mongodb
	span.LogKV("task", "delete user information from mongodb using email Id")
	err = info.Delete(span, email)
	zap.S().Info("info error ", err)
	if err != nil {
		span.LogKV("task", "send final output")
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "account successfully deleted",
	})
}
