package api

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/ginjwt"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/methods"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/signin"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/signup"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"net/http"
	"strings"
)

// User is a structure for binding data from signup or login request
type User struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignupUser is an api handler for creating account in postgres
func SignupUser(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "signup user")

	span.LogKV("task", "binding body data")
	// bind body data with variable
	var data User
	if err := c.BindJSON(&data); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": "Please pass email and password"})
		return
	}
	//saving data in account structure
	acs := dbtypes.User{}
	acs.Email = strings.ToLower(data.Email)
	acs.Role = "user"

	span.LogKV("task", "validating password")
	//validation check on password
	flag := methods.ValidatePassword(data.Password)
	if flag == 1 {
		c.JSON(400, gin.H{"error": true, "message": "Minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character."})
		return
	}
	span.LogKV("task", "save hash password")
	// save hash password instead of normal password
	acs.Password = methods.HashForNewPassword(data.Password)

	span.LogKV("task", "signup user")
	code, err := signup.Signup(acs)
	if err != nil {
		zap.S().Error(err)
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogKV("task", "send final output")
	c.JSON(200, gin.H{
		"error":   false,
		"message": "User account created successfully",
	})
}

// LoginUser is an api handler for checking credentials in postgres and generate jwt token
func LoginUser(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "login a user")

	span.LogKV("task", "binding body data")
	// bind body data with variable
	var data User
	if err := c.BindJSON(&data); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": "Please pass email and password"})
		return
	}

	span.LogKV("task", "check login credentials")
	login, msg, acc := signin.Signin(strings.ToLower(data.Email), data.Password)

	if !login {
		span.LogKV("task", "send final output")
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": msg})
		return
	} else {
		span.LogKV("task", "generate jwt token")
		mapd := ginjwt.GinJwtToken(acc)
		mapd["email"] = acc.Email
		mapd["role"] = acc.Role

		span.LogKV("task", "send final output")
		// passing user details and token details in a map with 200 status code
		c.JSON(http.StatusOK, mapd)
		return
	}

}
