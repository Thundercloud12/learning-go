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

func (apiConfig *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request)  {
	
	type parameters struct{
		Name string `name`
	}

	decoder := json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON:",err))
		return
	}

  	user,err:=apiConfig.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Couldnt create user:",err))
		return
	}
	
	respondWithJSON(w,200,databaseUserToUser(user))
}

func (apiConfig *apiConfig)handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User)  {
	respondWithJSON(w,200,databaseUserToUser(user))
}