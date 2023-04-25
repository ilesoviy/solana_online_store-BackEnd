package utils

import (
	"errors"
	"net/http"
	"time"

	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoInject the middleware base
var cli *mongo.Client
var _dbURI string
var _ctxDBinjectKey interface{}
var _ctxDBnamekey interface{}
var _ctxRawCliInjectKey interface{}

// Init sets all connection parameters
func DBInit(dbURI string, ctxDBinjectKey interface{}, ctxDBnamekey interface{}, ctxRawCliInjectKey interface{}) (context.CancelFunc, error) {
	if cli != nil {
		return nil, errors.New("already inited")
	}
	if ctxDBnamekey == nil || ctxRawCliInjectKey == nil || ctxDBinjectKey == nil {
		return nil, errors.New("invalid ctx keys specified")
	}
	_dbURI = dbURI
	_ctxDBnamekey = ctxDBnamekey
	_ctxDBinjectKey = ctxDBinjectKey
	_ctxRawCliInjectKey = ctxRawCliInjectKey
	//connect
	var err error
	cli, err = mongo.NewClient(options.Client().ApplyURI(_dbURI))
	if err != nil {
		return nil, err
	}
	ctx, ctxCancel := context.WithTimeout(context.Background(), 8*time.Second)
	if err := cli.Connect(ctx); err != nil {
		ctxCancel()
		return nil, err
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		ctxCancel()
		return nil, err
	}
	logrus.Infoln("DB connected")
	return ctxCancel, nil
}

// InjectClientMiddleware function, injects the raw DB client without a specific DB
func InjectDBClientMiddleware(c context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), _ctxRawCliInjectKey, cli))
		next.ServeHTTP(w, r)
	})
}

// InjectDBMiddleware function, injects the DB taking the DB name from config and the current DB session from session
func InjectDBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbName := viper.GetString("db.name")
		if dbName == "" {
			logrus.Errorln(errors.New("empty DB name provided by context to DBmiddleware"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		thedb := cli.Database(viper.GetString("db.name"))
		r = r.WithContext(context.WithValue(r.Context(), _ctxDBnamekey, dbName))
		r = r.WithContext(context.WithValue(r.Context(), _ctxDBinjectKey, thedb))
		next.ServeHTTP(w, r)
	})
}
