# Requirements => #
1. Mongo database to save user information
2. Redis database
3. Jaeger tracer image

# Environment Variables To be set to run=> #

```
// expose address
export MICROSERVICE2_HTTP_PORT="7002"

//private key for using authenticated rpc over different services to transfer data using token
export PRIVATE_KEY="abc"

// mysql environment variables
export MONGO_DB_NAME="grpc_example"
export MONGO_DB_HOST="locahost:27017"
export MONGO_DB_USER="username"
export MONGO_DB_PASS="password"
export MONGO_AUTH_DB_NAME="admin"

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
iii.  Run the app or binary -> command -> `$GOPATH/bin/micro-service2 --conf=environment <set flags for logs setting>`
```

### 3. Configuration using TOML file ###

```
i.    Create a configuration toml file by taking reference from example.toml file
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/micro-service2 --conf=toml --file=<path of toml file>  <set flags for logs setting>`
```

`Note :- for any help regarding flags, run this command '$GOPATH/bin/micro-service2 --help'`

## How to make changes in the project ##

1. Install golang and setup path (reference -> https://golang.org/doc/install)
2. Clone repo and save it in `$GOPATH/src/git.xenonstack.com/util/`
3. Install [go-dep](https://github.com/golang/dep) for updating any dependencies
4. Update dependencies -> command -> `$dep ensure -update` (mainly for changes in-house libraries)
5. Update the code if any.
6. Test the changes locally.
7. Push the changes and also push the changes of Gopkg.lock file and vendor folder also.
