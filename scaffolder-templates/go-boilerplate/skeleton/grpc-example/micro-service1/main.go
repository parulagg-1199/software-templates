package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"git.xenonstack.com/lib/golang-boilerplate/logger"
	"git.xenonstack.com/lib/golang-boilerplate/tracing"

	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/config"
	pb "git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/pb"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/src/dbtypes"
	"git.xenonstack.com/util/golang-boilerplate/grpc-example/micro-service1/src/services"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	ot "github.com/opentracing/opentracing-go"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	// to save grpc server otions
	var opts []grpc.ServerOption
	// to save grpc unary server interceptors
	var unarys []grpc.UnaryServerInterceptor
	// to save grpc stream server interceptors
	var streams []grpc.StreamServerInterceptor

	//enable logger
	zap.S().Info("enable logger")
	// for replacing grpc default logger and also global logger also
	unary, stream, err := logger.AddLogging(-1, *logType, *logEnv)
	if err != nil {
		// if any error
		zap.S().Error(err)
	} else {
		// save in slice of unary server acceptor and stream server acceptor
		unarys = append(unarys, unary)
		streams = append(streams, stream)
	}

	//enable tracer
	zap.S().Info("enable tracer")

	// fetch deault configuration setting of tracer
	confTrace := tracing.New()
	// change setting according to service
	confTrace.ServiceName = "grpc-user-service"
	confTrace.JaegerHost = config.Conf.Jaeger.Host
	confTrace.JeagerPort = config.Conf.Jaeger.Port
	// for adding tracer in grpc
	unary, stream, tracer, closer := confTrace.GrpcServerOptions()
	// close when all tracer operations completed
	defer closer.Close()
	// set this tracer as global tracer
	ot.SetGlobalTracer(tracer)
	// save in slice of unary server acceptor and stream server acceptor
	unarys = append(unarys, unary)
	streams = append(streams, stream)

	//initialise mysql
	go dbtypes.CreateMySQLDBTablesIfNotExists()

	zap.S().Info("intialize grpc server options")
	// now add all interceptors in grpc options
	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unarys...)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streams...)))
	zap.S().Info(opts)

	//start grpc server
	server := grpc.NewServer(opts...)
	serverInterface := &services.Server{
		Service2Client: services.GetS2Client(),
	}
	//register service with grpc
	pb.RegisterUserDetailsServer(server, serverInterface)

	if config.Conf.Service.Environment == "grpc" {

		// create tcp connection
		zap.S().Info("create tcp connection")
		lis, err := net.Listen("tcp", ":"+config.Conf.Service.Port)
		if err != nil {
			zap.S().Error("failed to initialize TCP listen: %v", err)
			return
		}

		//close tcp connection
		defer lis.Close()
		//start listening at tcp connection
		err = server.Serve(lis)
		if err != nil {
			zap.S().Error("failed to serve grpc: %v", err)
			return
		}
		os.Exit(1)
	}

	//conver grpc server to web grpc server
	wrappedGrpc := grpcweb.WrapServer(server)
	// log.Println(wrappedGrpc.IsAcceptableGrpcCorsRequest())
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedGrpc.ServeHTTP(resp, req)
	}

	// configure http server
	httpServer := http.Server{
		Addr:    fmt.Sprint(":", config.Conf.Service.Port),
		Handler: http.HandlerFunc(handler),
	}

	zap.S().Info("Starting server. http port: ", config.Conf.Service.Port)

	//start listening http server
	err = httpServer.ListenAndServe()
	if err != nil {
		zap.S().Error("failed to serve grpc: ", err)
		return
	}
}
