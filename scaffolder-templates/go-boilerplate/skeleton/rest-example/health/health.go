package health

import (
	"git.xenonstack.com/util/golang-boilerplate/rest-example/config"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"errors"
)

//Healthz to check health of service -> check each resource health
func Healthz(parentspan opentracing.Span) error {
	span := opentracing.StartSpan("health function", opentracing.ChildOf(parentspan.Context()))
	defer span.Finish()
	span.LogKV("check", "postgres")
	//check health of postgres
	ok := config.Conf.Postgres.PostgresHealthz()
	zap.S().Info("postgre", ok)
	if !ok {
		span.LogKV("return ", "postgres error")
		return errors.New("PostgreSQL is not working")
	}
	span.LogKV("check", "MongoDB")
	//check health of mongo
	ok = config.Conf.Mongo.MongoHealthz()
	zap.S().Info("mongo", ok)
	if !ok {
		span.LogKV("return ", "MongoDB error")
		return errors.New("MongoDB is not working")
	}
	span.LogKV("check", "postgres")
	//check health of mysql
	ok = config.Conf.Mysql.MysqlHealthz()
	zap.S().Info("mysql", ok)
	if !ok {
		span.LogKV("return ", "MySQL error")
		return errors.New("MySQL is not working")
	}
	span.LogKV("check", "postgres")
	//check health of redis
	ok = config.Conf.Redis.Healthz()
	zap.S().Info("redis ", ok)
	if !ok {
		span.LogKV("return ", "Redis error")
		return errors.New("Redis is not working")
	}
	span.LogKV("return ", "nil error")
	return nil
}
