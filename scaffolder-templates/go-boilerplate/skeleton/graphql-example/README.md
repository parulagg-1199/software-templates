# Requirements => #
1. Postgresql database to save accounts
2. Jaeger tracer image

# Environment Variables To be set to run=> #

```
// expose port
export GRAPHQL_EXAMPLE_HTTP_PORT="8001"

//admin account environment variables
export GRAPHQL_EXAMPLE_ADMIN_EMAIL="admin@boilerplate.com"
export GRAPHQL_EXAMPLE_ADMIN_PASS="boiler"

//private key for using authenticated apis over different services to transfer data using token
export GRAPHQL_EXAMPLE_PRIVATE_KEY="abc"

// postgresql environment variables
export POSTGRESQL_DB_HOST="localhost"
export POSTGRESQL_DB_PORT="5432"
export POSTGRESQL_DB_USER="postgres"
export POSTGRESQL_DB_PASS="1234"
export POSTGRESQL_DB_NAME="graphql_example"
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
iii.  Run the app or binary -> command -> `$GOPATH/bin/server --conf=environment <set flags for logs setting>`
```

### 3. Configuration using TOML file ###

```
i.    Create a configuration toml file by taking reference from example.toml file
ii.   Build the app or binary -> command -> `$ go install`
iii.  Run the app or binary -> command -> `$GOPATH/bin/server --conf=toml --file=<path of toml file>  <set flags for logs setting>`
```

`Note :- for any help regarding flags, run this command '$GOPATH/bin/server --help'`


# How to make changes in the project => #

1. Install golang and setup path (reference -> https://golang.org/doc/install)
2. Clone repo and save it in `$GOPATH/src/git.xenonstack.com/lib/`
3. Install [go-dep](https://github.com/golang/dep) for updating any dependencies
4. Update dependencies -> command -> `$dep ensure -update` (mainly for changes in-house libraries)
5. Update the code if any.
6. Fulfill above defined requirements
6. Test the changes locally -> Run command `go run server/server.go`.
7. Push the changes and also push the changes of Gopkg.lock file and vendor folder also.

# Getting started with graphql => #

1.Define the schema in schema.graphql file using the [Graphql Schema Definition Language](http://graphql.org/learn/schema/)<br/>
2.Create the project skeleton -> `go run github.com/99designs/gqlgen init`
This has created an empty skeleton with all files you need:
* gqlgen.yml — The gqlgen config file, knobs for controlling the generated code.
* generated.go — The GraphQL execution runtime, the bulk of the generated code.
* models_gen.go — Generated models required to build the graph. Often you will override these with your own models. Still very useful for input types.
* resolver.go — This is where your application code lives. generated.go will call into this to get the data the user has requested.
* server/server.go — This is a minimal entry point that sets up an http.Handler to the generated GraphQL server.<br/>

3.Create the database models or new models and Next tell gqlgen to use this new struct by adding it to `gqlgen.yml` :

```
models:
  Todo:
    model: github.com/[username]/gqlgen-todos.Todo
```

4.Regenerate the file by running: `go run github.com/99designs/gqlgen`<br/>
5.Implement the Resolvers<br/>
6.run app by `go run server/server.go`<br/>
