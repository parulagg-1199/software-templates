package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/opengintracing"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ot "github.com/opentracing/opentracing-go"
	logs "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"

	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/health"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/routes"

	"git.xenonstack.com/lib/golang-boilerplate/logger"
	"git.xenonstack.com/lib/golang-boilerplate/tracing"

	"flag"
	"net/http"
	"os"
	"time"
)

func main() {
	// setup for reading flags for deciding whether to do configuration with toml or env variables
	// and environment and format type for logging
	conf := flag.String("conf", "environment", "set configuration from toml file or environment variables")
	file := flag.String("file", "", "set path of toml file")
	logType := flag.String("logType", "tab space", "set log type as tab space or json")
	logEnv := flag.String("logEnv", "production", "set log environment as development or production")
	flag.Parse()

	// flag parsing ofr configuration
	if *conf == "environment" {
		log.Println("environment")
		config.ConfigurationWithEnv()
	} else if *conf == "toml" {
		log.Println("toml")
		if *file == "" {
			log.Println("Please provide file path if you want to configure with toml file")
			flag.PrintDefaults()
			os.Exit(1)
		} else {
			config.ConfigurationWithToml(*file)
		}
	} else {
		log.Println("Please pass valid arguments, conf can be set as toml or environment")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// initialize logger
	err := logger.Init(-1, *logType, *logEnv)
	if err != nil {
		zap.S().Error(err)
	}

	// fetch deault configuration setting of tracer
	confTrace := tracing.New()
	// change setting according to service
	confTrace.ServiceName = "rest-example"
	confTrace.JaegerHost = config.Conf.Jaeger.Host
	confTrace.JeagerPort = config.Conf.Jaeger.Port

	// initialize tracer
	trace, closer := confTrace.InitJaeger()
	defer closer.Close()
	ot.SetGlobalTracer(trace)

	//initialise postgres
	go dbtypes.CreatePostgreDBTablesIfNotExists()
	//initialise mysql
	go dbtypes.CreateMySQLDBTablesIfNotExists()

	// initialize gin router
	router := gin.Default()
	//set zap logger as std logger
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))

	//allowing CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowMethods("DELETE")
	router.Use(cors.New(corsConfig))

	zap.S().Info("health endpoint")
	// api route to check health of serive
	router.GET("/health",
		opengintracing.NewSpan(ot.GlobalTracer(), "checking health of service"),
		healthz)
	// defined other routes
	routes.V1Routes(router)

	zap.S().Info(config.Conf.Service.Port)
	// Listen and Server in 0.0.0.0:8080
	router.Run(":" + config.Conf.Service.Port)
}

func healthz(c *gin.Context) {
	span, found := opengintracing.GetSpan(c)
	if found == false {
		zap.S().Warn("Span not found")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	span.SetTag("event", "check health")
	span.LogKV("info", "start health check")
	err := health.Healthz(span)
	if err != nil {
		span.LogFields(
			logs.String("info", "stop health check"),
			logs.String("error", err.Error()),
		)
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	span.LogFields(
		logs.String("info", "stop health check"),
		logs.String("error", "nil"),
	)
	c.JSON(200, gin.H{
		"error":   true,
		"message": "ok",
	})
	return
}
