# Back-end

## Start back-end
### Start database
- Install **docker** and **docker-compose plugin** if not installed with docker installation. In windows maybe only option is *docker desktop* (https://docs.docker.com/engine/install/)
- Go to *docker* folder, next run `docker-compose up`. Localhost db connection string is **mongodb://localhost:27017** and the default db name used in back-end is **sol_project**
- In order to visualize data in db and insert them, we can use MongoDBCompass (https://www.mongodb.com/try/download/compass)
### Start application
- Install golang (https://go.dev/doc/install)
- run `go run main.go`
- api endpoints are in *v0/* prefix (ex. http://localhost:4000/v0/items/populars)

### Configuration file
In order to run different config file, we must run `go run main.go -cfg {filename}`, without exstension. The default config file is `config_debug.json`
