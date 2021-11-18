# Requirements => #
1. Mysql database to save user details
2. Redis database
4. Jaeger tracer image

# Environment Variables To be set to run=> #

```
// expose port
export MICROSERVICE1_HTTP_PORT="7000"

// 2nd service Address
export MICROSERVICE2_HOST_ADDRESS="localhost:7002"

//private key for using authenticated rpc over different services to transfer data using token
export PRIVATE_KEY="abc"

// this defined whether to deploy server as grpc or http
export ENVIRONMENT="grpc"

// mysql environment variables
export MYSQL_DB_HOST="locahost"
export MYSQL_DB_PORT="3306"
export MYSQL_DB_USER="root"
export MYSQL_DB_PASS=""
export MYSQL_DB_NAME="rest_example"

// redis environment variables
export REDIS_HOST="localhost"
export REDIS_PORT="6379"
export REDIS_PASS=""
export REDIS_DB=0

//jaeger configuration
export JAEGER_AGENT_HOST="localhost"
export JAEGER_AGENT_PORT="6831"
```

## How to run the app ##

### 1. Flags Information ###

```
i.    conf    ->  set configuration from toml file or environment variables (default "environment")
ii.   file    ->  set path of toml file if configuration if set to toml
iii.  logEnv  ->  set log environment as development or production (default "production") change in output result
iv.   logType ->  set log format as tab space or json (default "tab space")
```

### 2. Configuration using environment variables ###

```
i.    Export above all environment variables
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/micro-service1 --conf=environment <set flags for logs setting>`
```

### 3. Configuration using TOML file ###

```
i.    Create a configuration toml file by taking reference from example.toml file
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/micro-service1 --conf=toml --file=<path of toml file>  <set flags for logs setting>`
```

`Note :- for any help regarding flags, run this command '$GOPATH/bin/micro-service1 --help'`

## How to make changes in the project ##

1. Install golang and setup path (reference -> https://golang.org/doc/install)
2. Clone repo and save it in `$GOPATH/src/git.xenonstack.com/util/`
3. Install [go-dep](https://github.com/golang/dep) for updating any dependencies
4. Update dependencies -> command -> `$dep ensure -update` (mainly for changes in-house libraries)
5. Update the code if any.
6. Test the changes locally.
7. Push the changes and also push the changes of Gopkg.lock file and vendor folder also.

# Getting Started with grpc server =>#
1. Create protobuf file with `.proto` extension according to your service requirements.
2. Install the protoc compiler and protoc plugin for go using this command `go get -u github.com/golang/protobuf/protoc-gen-go`
3. Compile the proto file means generating code for your choosen language (GoLang). Command for compilation `protoc --proto_path=<proto folder> --go_out=plugins=grpc:<path where you want your generated file> <path/file_name.proto>`
4. Make the grpc server file using above generated file
5. Then run that file to run the server
