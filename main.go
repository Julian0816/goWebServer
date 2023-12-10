package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Julian0816/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load() // Loads the environment variables

	portString := os.Getenv("PORT")
	if (portString == "") {
		log.Fatal("PORT is not found in the environment") // Exit the program immediately with error code 1 and a message
	}

	dbURL := os.Getenv("DB_URL")
	if (dbURL == "") {
		log.Fatal("DB_URL is not found in the environment") // Exit the program immediately with error code 1 and a message
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err) //
	}
	
	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"https://*", "http://*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"*"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: false,
    MaxAge:           300,
   }))

   // The full path will be /v1/healthz (Good practice to check the health of the server)
   v1Router := chi.NewRouter()
   v1Router.Get("/healthz", handlerReadiness) // Connect the handlerReadiness to the "/healthz" path
   v1Router.Get("/err", handlerErr)
   v1Router.Post("/users", apiCfg.handlerCreateUser) // Create User handler

   router.Mount("/v1", v1Router) // nest a v1Router under the v1 path // This is good practive in case you need to make a v2 route

 	  

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	
	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port: ", portString)
}