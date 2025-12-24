package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Thundercloud12/go-projects/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

func main(){


	godotenv.Load(".env")
	portString:= os.Getenv("PORT")

	if portString==""{
		log.Fatal("PORT not found in the environment")
	}
	DBString:= os.Getenv("DB_URL")

	if DBString==""{
		log.Fatal("DB not found in the environment")
	}

	conn,err:=sql.Open("postgres",DBString)
	if err !=nil{
		fmt.Printf("Error aagaya %v",err)
	}
	db:=database.New(conn)

	apiCfg:=apiConfig{
		DB:db,
	}

	go startScraping(db,10,time.Minute)

	router:=chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
        ExposedHeaders: []string{"*"},
        AllowCredentials: false,
        MaxAge: 300,
    }))

	v1Router:=chi.NewRouter()

	v1Router.Get("/ready",handlerReadiness)
	v1Router.Get("/err",handlerErr)
	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feed",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Post("/subscribeFeeds",apiCfg.middlewareAuth(apiCfg.handlersubscribeFeeds))
	v1Router.Get("/allfeeds",apiCfg.handlerGetFeeds)
	v1Router.Get("/allfeedforuser",apiCfg.middlewareAuth(apiCfg.handlergetFeedFollower))
	v1Router.Post("/deletefeedfollower",apiCfg.middlewareAuth(apiCfg.handledeletefollower))
	router.Mount("/v1",v1Router)
	srv:=&http.Server{
		Handler: router,
		Addr: ":"+portString,
	}

	

	log.Printf("Server started on port %v", portString)

	srv.ListenAndServe()

	fmt.Println("Port:",portString)
}