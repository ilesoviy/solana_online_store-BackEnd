package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/SolProj3ct/Back-end/integrations/algolia"
	"github.com/SolProj3ct/Back-end/routes"
	"github.com/SolProj3ct/Back-end/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	cfgFile := flag.String("cfg", "config_debug", "Specify the viper config file to be used.")
	flag.Parse()
	// load config file
	initConfig(*cfgFile)
	// connect to db
	ctxCancel, e := utils.DBInit(viper.GetString("db.uri"), utils.DB, utils.DBName, utils.DBCli)
	if e != nil {
		panic(e)
	}
	defer ctxCancel()
	// initiate magic service
	/*if err := magic.Init(viper.GetString("integrations.magic.apiSecret")); err != nil {
		panic(err)
	}*/
	// initiate algolia service
	if err := algolia.Init(viper.GetString("integrations.algolia.appID"), viper.GetString("integrations.algolia.apiSecret")); err != nil {
		panic(err)
	}
	// create api router
	router := setupApiRoutes()
	// create interrupt signal listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	// start webserver
	fmt.Println(banner)
	srv := &http.Server{Addr: ":4000", Handler: router}
	go func() {
		fmt.Println("SERVING BACKEND ON PORT 4000")
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-quit
	logrus.Infoln("Shutting down")
}

func initConfig(cfgFile string) {
	viper.SetConfigType("json")
	viper.SetDefault("logger.level", "3")
	viper.SetConfigName(cfgFile) // name of config file (without extension)
	viper.AddConfigPath(".")     // path to look for the config file in
	if cfgErr := viper.ReadInConfig(); cfgErr != nil {
		panic(cfgErr)
	}
	logrus.Infof("Static config loaded from file \"%s\"\n", cfgFile)
}

func setupApiRoutes() *mux.Router {
	r := mux.NewRouter()
	authApi := r.PathPrefix("/v0").Subrouter()
	unAuthApi := r.PathPrefix("/auth").Subrouter()

	// db middleware
	r.Use(utils.InjectDBMiddleware)

	// AUTH API
	// jwt verification middleware
	authApi.Use(utils.InjectUserMiddleware)
	// items
	authApi.HandleFunc("/items/brands", routes.GetItemBrands).Methods("GET")
	authApi.HandleFunc("/items/recommended", routes.GetRecommendedItems).Methods("GET")
	authApi.HandleFunc("/items/populars", routes.GetPopularItems).Methods("GET")
	authApi.HandleFunc("/items", routes.NewItem).Methods("POST")
	// search
	authApi.HandleFunc("/search", routes.GlobalSearch).Methods("GET")
	// users
	authApi.HandleFunc("/users", routes.NewUser).Methods("POST")

	// UNAUTH
	unAuthApi.HandleFunc("/signin", routes.SignIn).Methods("POST")
	unAuthApi.HandleFunc("/login", routes.LogIn).Methods("POST")
	return r
}

const banner = `
  ____        _ ____            _           _   
 / ___|  ___ | |  _ \ _ __ ___ (_) ___  ___| |_ 
 \___ \ / _ \| | |_) | '__/ _ \| |/ _ \/ __| __|
  ___) | (_) | |  __/| | | (_) | |  __/ (__| |_ 
 |____/ \___/|_|_|   |_|  \___// |\___|\___|\__|
                             |__/             `
