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
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Link struct {
	ID  string
	URL string
}

var LinkMap = map[string]*Link{
	"example": {ID: "example", URL: "https://example.com"},
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

	LinkMap[id] = &Link{ID: id, URL: url}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (apiconfig *apiconfig)RedirectHandler(c *echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.String(http.StatusBadRequest, "ID missing")
	}

	link, found := LinkMap[id]
	if !found {
		return c.String(http.StatusNotFound, "Aee jasti shanpati nako karu")
	}

	return c.Redirect(http.StatusMovedPermanently, link.URL)
}

func main() {


	godotenv.Load(".env")
	portString:=os.Getenv("PORT")

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

	apiCfg:=apiconfig{
		DB:db,
	}



	e := echo.New()

	// Logging middleware
	e.Use(middleware.RequestLogger())

	component := landing()

	// Home page (templ)
	e.GET("/", echo.WrapHandler(templ.Handler(component)))

	// Submit short link
	e.POST("/submit", apiCfg.SubmitHandler)

	// Redirect handler (keep last to avoid route hijacking)
	e.GET("/:id", apiCfg.RedirectHandler)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
