package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"url-short/internal/database"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	_ "github.com/lib/pq"

)

type Link struct {
	ID  uuid.UUID
	URL string
}


type apiconfig struct{
	DB *database.Queries
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Seed once globally
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	var result []byte
	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}
	return string(result)
}

func (apiconfig *apiconfig)SubmitHandler(c *echo.Context) error {
	url := c.FormValue("url")

	if url == "" {
		return c.String(http.StatusBadRequest, "Link khaali ka pathavtoys re")
	}

	// Add https:// if missing
	if !(len(url) >= 4 && (url[:4] == "http" || url[:5] == "https")) {
		url = "https://" + url
	}

	id := generateRandomString(8)

	// LinkMap[id] = &Link{ID: id, URL: url}

	err:= apiconfig.DB.CreateMapping(c.Request().Context(),database.CreateMappingParams{
		ID: uuid.New(),
		OriginalUrl: url,
		ConvertedUrl: id,
	})

	if err!=nil{
		
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (apiconfig *apiconfig)RedirectHandler(c *echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.String(http.StatusBadRequest, "ID missing")
	}

	link, err:= apiconfig.DB.GetConverted(c.Request().Context(),id)
	if err!=nil{
		
		return c.String(http.StatusInternalServerError, "Correct link pathav")
	}

	return c.Redirect(http.StatusMovedPermanently, link.OriginalUrl)
}

func (apiconfig *apiconfig)HomeHandler() echo.HandlerFunc {

	  return func(c *echo.Context) error {
		links,err:= apiconfig.DB.GetAllMappings(c.Request().Context())

		if err != nil{
			return c.String(500,err.Error())
		}

		component:=landing(links)

		return echo.WrapHandler(templ.Handler(component))(c)
	  }
	
}

func main() {


	godotenv.Load(".env")
	// portString:=os.Getenv("PORT")

	// if portString==""{
	// 	log.Fatal("PORT not found in the environment")
	// }

	DBString:= os.Getenv("DB_URL")

	if DBString==""{
		log.Fatal("DB not found in the environment")
	}

	conn,err:=sql.Open("postgres",DBString)
	if err !=nil{
		fmt.Printf("Error aagaya %v",err)
	}
	db:=database.New(conn)

	apiCfg:=apiconfig{
		DB:db,
	}



	e := echo.New()

	// Logging middleware
	e.Use(middleware.RequestLogger())

	

	// Home page (templ)
	e.GET("/", apiCfg.HomeHandler())

	// Submit short link
	e.POST("/submit", apiCfg.SubmitHandler)

	// Redirect handler (keep last to avoid route hijacking)
	e.GET("/:id", apiCfg.RedirectHandler)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
