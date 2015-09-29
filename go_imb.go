package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
)

type GoImbConfiguration struct {
		Port int
		Webapp string
		MongodbUrl string
		MongodbDatabaseName string
}

func readConfiguration() GoImbConfiguration {
	confFilePath := os.Getenv("GOIMBCONF")
	if len(confFilePath) == 0  {
		confFilePath = "src/github.com/martianov/go-imb/go-imb.conf.default"
	}
	configFileData, readErr := ioutil.ReadFile(confFilePath)
	if readErr != nil {
        panic(fmt.Errorf("Failed to read configuration file %s: %v", confFilePath, readErr))
    }

	configuration := GoImbConfiguration{}
	parseErr := yaml.Unmarshal([]byte(configFileData), &configuration)
	if parseErr != nil {
		panic(fmt.Errorf("Failed to parse configuration file %s: %v", confFilePath, parseErr))
	}

	fmt.Printf("Configuration file: %s\n", confFilePath)

    return configuration
}

func main() {
	configuration := readConfiguration()
	dbSession := DBConnect(configuration.MongodbUrl)
	DBEnsureIndicesAndDefaults(dbSession, configuration.MongodbDatabaseName)

	// handle all requests by serving a file of the same name
	fs := http.Dir(configuration.Webapp)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()

	router.Handle("/", http.RedirectHandler("/webapp/index.html", 302))
	router.PathPrefix("/webapp").Handler(http.StripPrefix("/webapp", fileHandler))

	authRouterBase := mux.NewRouter();
	router.PathPrefix("/auth").Handler(negroni.New(DBMiddleware(dbSession, configuration.MongodbDatabaseName), negroni.Wrap(authRouterBase)))
	authRouter := authRouterBase.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", Login).Methods("POST")

	apiRouterBase := mux.NewRouter();
	router.PathPrefix("/api").Handler(negroni.New(DBMiddleware(dbSession, configuration.MongodbDatabaseName), JWTMiddleware(), negroni.Wrap(apiRouterBase)))
	apiRouter := apiRouterBase.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/me", Me).Methods("GET")

	http.ListenAndServe(fmt.Sprintf(":%v", configuration.Port), router);
}