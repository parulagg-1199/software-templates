package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	graphql_example "git.xenonstack.com/util/golang-boilerplate/graphql-example"
	"git.xenonstack.com/util/golang-boilerplate/graphql-example/config"
	"git.xenonstack.com/util/golang-boilerplate/graphql-example/src/auth"
	"git.xenonstack.com/util/golang-boilerplate/graphql-example/src/postgredb"

	"git.xenonstack.com/lib/golang-boilerplate/logger"
	"git.xenonstack.com/lib/golang-boilerplate/tracing"

	"github.com/99designs/gqlgen/handler"
	ot "github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"go.uber.org/zap"
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

	//initialise postgres database
	postgredb.CreatePostgreDBTablesIfNotExists()

	// connect with postgres
	db, err := postgredb.Connect()
	if err != nil {
		zap.S().Error(err)
		os.Exit(1)
	}

	// initialize http multiplexer
	mux := http.NewServeMux()
	// api endpoint for GraphQL PLayGround
	mux.Handle("/", handler.Playground("GraphQL playground", "/query"))

	// api endpoint for query execution
	mux.Handle("/query", auth.Middleware(handler.GraphQL(graphql_example.NewExecutableSchema(graphql_example.NewRootResolvers(db)))))

	// enable cors
	handler := cors.AllowAll().Handler(mux)
	zap.S().Infof("connect to http://localhost:%v/ for GraphQL playground", config.Conf.Service.Port)
	http.ListenAndServe(":"+config.Conf.Service.Port, handler)
}
