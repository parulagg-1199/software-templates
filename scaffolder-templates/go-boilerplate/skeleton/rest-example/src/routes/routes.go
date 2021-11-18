package routes

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/api"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/ginjwt"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// V1Routes is a method in which all the service endpoints are defined
func V1Routes(router *gin.Engine) {
	v1 := router.Group("/v1")

	zap.S().Info("un-protected apis")
	// endpoints related to account creation and login with postgre
	v1.POST("/signup", opengintracing.NewSpan(ot.GlobalTracer(), "signup user"), api.SignupUser)
	v1.POST("/login", opengintracing.NewSpan(ot.GlobalTracer(), "login user"), api.LoginUser)
	// v1.DELETE("/user/:id", opengintracing.NewSpan(ot.GlobalTracer(), "delete user"), api.DeleteUser)

	//setting up middleware for protected apis
	authMiddleware := ginjwt.MwInitializer()

	zap.S().Info("protected apis")
	//Protected apis
	v1.Use(opengintracing.NewSpan(ot.GlobalTracer(), "gin jwt middleware"), authMiddleware.MiddlewareFunc())
	{
		// api endpoint related to mysql
		v1.POST("/userDetail", opengintracing.NewSpan(ot.GlobalTracer(), "save user details"), api.AddUserDetails)
		v1.GET("/userDetail", opengintracing.NewSpan(ot.GlobalTracer(), "get user details"), api.GetUserDetails)

		// api endpoint related to mongodb
		v1.POST("/userInfo", opengintracing.NewSpan(ot.GlobalTracer(), "save user information"), api.SaveInformation)
		v1.GET("/userInfo", opengintracing.NewSpan(ot.GlobalTracer(), "fetch user information"), api.FetchInformation)

		// api endpoints related to redis
		v1.POST("/key", opengintracing.NewSpan(ot.GlobalTracer(), "save key in redis"), api.SaveData)
		v1.GET("/key", opengintracing.NewSpan(ot.GlobalTracer(), "check key in redis"), api.CheckKey)
		v1.DELETE("/key", opengintracing.NewSpan(ot.GlobalTracer(), "delete key from redis"), api.DeleteKey)

		// add custom middleware for checking person is admin or not
		v1.Use(opengintracing.NewSpan(ot.GlobalTracer(), "custom middleware"), api.CheckAdmin)
		{
			//delete user api by admin only
			v1.DELETE("/deleteAccount", opengintracing.NewSpan(ot.GlobalTracer(), "delete account"), api.DeleteAccount)
		}
	}
}
