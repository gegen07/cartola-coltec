package main

import (
	"github.com/gegen07/cartola-university/application"
	persistence "github.com/gegen07/cartola-university/infrastructure/persistance"
	"github.com/gegen07/cartola-university/interfaces"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env")
	}
}

func main() {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbpass := os.Getenv("DB_PASSWORD")
	dbuser := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	dbdriver := os.Getenv("DB_DRIVER")

	repositories, err := persistence.NewRepositories(dbdriver, dbuser, dbpass, dbport, dbhost, dbname)

	if err != nil {
		panic(err)
	}

	positionApplication := application.NewPositionApplication(repositories.Position, repositories.ScoutPosition, time.Second*5)
	scoutApplication := application.NewScoutApplication(repositories.Scout, repositories.ScoutPosition, time.Second*5)

	position := interfaces.NewPositionHandler(positionApplication)
	formations := interfaces.NewFormationHandler(repositories.Formation)
	scouts := interfaces.NewScoutHandler(scoutApplication)

	router := mux.NewRouter()

	//formation routes
	router.Handle("/formation", interfaces.RootHandler(formations.Insert)).Methods(http.MethodPost)
	router.Handle("/formation", interfaces.RootHandler(formations.GetAllFormations)).Methods(http.MethodGet).Queries("page", "{page}")
	router.Handle("/formation/{id}", interfaces.RootHandler(formations.GetFormationById)).Methods(http.MethodGet)
	router.Handle("/formation/{id}", interfaces.RootHandler(formations.Update)).Methods(http.MethodPut)
	router.Handle("/formation/{id}", interfaces.RootHandler(formations.Delete)).Methods(http.MethodDelete)

	//scout routes
	router.Handle("/scout", interfaces.RootHandler(scouts.Insert)).Methods(http.MethodPost)
	router.Handle("/scout", interfaces.RootHandler(scouts.GetAllScouts)).Methods(http.MethodGet).Queries("page", "{page}")
	router.Handle("/scout/{id}", interfaces.RootHandler(scouts.GetScoutByID)).Methods(http.MethodGet)
	router.Handle("/scout/{id}", interfaces.RootHandler(scouts.Update)).Methods(http.MethodPut)
	router.Handle("/scout/{id}", interfaces.RootHandler(scouts.Delete)).Methods(http.MethodDelete)

	//position routes
	router.Handle("/position", interfaces.RootHandler(position.Insert)).Methods(http.MethodPost)
	router.Handle("/position", interfaces.RootHandler(position.GetAll)).Methods(http.MethodGet).Queries("page", "{page}")
	router.Handle("/position/{id}", interfaces.RootHandler(position.GetByID)).Methods(http.MethodGet)
	router.Handle("/position/{id}", interfaces.RootHandler(position.Update)).Methods(http.MethodPut)
	router.Handle("/position/{id}", interfaces.RootHandler(position.Delete)).Methods(http.MethodDelete)

	logrus.Infof("Starting server ...")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(router)))
}
