package api

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CheckAdmin is a middleware function use to check person is admin or user
func CheckAdmin(c *gin.Context) {
	// fetch opentracing span from context
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(500)
		return
	}
	span.SetTag("event", "check person is user or not")

	span.LogKV("task", "fetch claims and check person is user or admin")
	// extract claims from jwt token
	claims := jwt.ExtractClaims(c)
	span.LogKV("task", "send final output")
	if claims["sys_role"].(string) != "admin" {
		// when person is not admin
		c.Abort()
		c.JSON(401, gin.H{"error": true, "message": "You are not allowed to perform these actions"})
		return
	}

	c.Next()
}
