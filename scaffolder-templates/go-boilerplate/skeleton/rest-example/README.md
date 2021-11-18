# Boilerplate Rest Example #

Rest Example is an example for how to create REST APIs in **Golang** with proper file structure. Other features included:-
* How to add middleware in routes
* Configuration management using toml file
* How to use in-house libraries(private packages) and how to update in vendor and Gopkg.toml file
* Example of each database queries -> mongo, postgresql, mysql and redis database
* JWT authentication with go-gin
* Each component health check
* Example of opentracing with go-gin api's
* Proper code-commenting according to golang standards

## Requirements => ##
1. Postgresql database to save accounts
2. Mysql database to save user details
3. Mongo database to save user informations
4. Redis to save key-value data
5. Jaeger tracer image

## Environment Variables To be set to run=> ##

```
// expose port
export REST_EXAMPLE_HTTP_PORT="8001"

//admin account environment variables
export REST_EXAMPLE_ADMIN_EMAIL="admin@boilerplate.com"
export REST_EXAMPLE_ADMIN_PASS="boiler"

//private key for using authenticated apis over different services to transfer data using token
export REST_EXAMPLE_PRIVATE_KEY="abc"

// postgresql environment variables
export POSTGRESQL_DB_HOST="localhost"
export POSTGRESQL_DB_PORT="5432"
export POSTGRESQL_DB_USER="postgres"
export POSTGRESQL_DB_PASS="1234"
export POSTGRESQL_DB_NAME="rest_example"

// mysql environment variables
export MYSQL_DB_HOST="locahost"
export MYSQL_DB_PORT="3306"
export MYSQL_DB_USER="root"
export MYSQL_DB_PASS=""
export MYSQL_DB_NAME="rest_example"

// mongodb environment variables
export MONGO_AUTH_DB_NAME="admin"
export MONGO_DB_NAME="rest_example"
export MONGO_DB_HOST="localhost:27017"
export MONGO_DB_USER="username"
export MONGO_DB_PASS="password"

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
iii.  Run the app or binary -> command -> `$GOPATH/bin/rest-example --conf=environment <set flags for logs setting>`
```

### 3. Configuration using TOML file ###

```
i.    Create a configuration toml file by taking reference from example.toml file
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/rest-example --conf=toml --file=<path of toml file>  <set flags for logs setting>`
```

`Note :- for any help regarding flags, run this command '$GOPATH/bin/rest-example --help'`

## How to make changes in the project ##

1. Clone the repo
2. Install [go-dep](https://github.com/golang/dep) for updating any dependencies
3. Update dependencies -> command -> `$dep ensure -update` (mainly for changes in-house libraries)
4. Update the code if any.
5. Test the changes locally.
6. Push the changes and also push the changes of Gopkg.lock file and vendor folder also.
