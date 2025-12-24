package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_"github.com/Thundercloud12/go-projects/internal/auth"
	"github.com/Thundercloud12/go-projects/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig)handlerCreateFeed(w http.ResponseWriter, r *http.Request,user database.User)  {
	
	type parameters struct{
		Name string `json":name"`
		Url string  `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:",err))
		return
	}

  	feed,err:=apiConfig.DB.CreatrFeed(r.Context(),database.CreatrFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldnt create user:",err))
		return
	}
	
	respondWithJSON(w,200,databaseFeedtoFeed(feed))
}

func (apiConfig *apiConfig)handlerGetFeeds(w http.ResponseWriter, r *http.Request)  {
	
  	feed,err:=apiConfig.DB.FetchFeeds(r.Context())
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error fetching feeds",err))
		return
	}

	respondWithJSON(w,200,databaseFeedstoFeeds(feed))
}


