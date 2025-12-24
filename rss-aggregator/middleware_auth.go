package main

import (
	"fmt"
	"net/http"

	"github.com/Thundercloud12/go-projects/internal/auth"
	"github.com/Thundercloud12/go-projects/internal/database"
)


type authedHandler func(http.ResponseWriter,*http.Request, database.User)

func (apiConfig *apiConfig) middlewareAuth(handler authedHandler)  http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request) {
		apikey,err:=auth.GetAPIey(r.Header)
		if err!=nil{
			respondWithError(w,403,fmt.Sprintf("Auth error: %v",err))
			return 
		}

		user,err:=apiConfig.DB.GetUserByAPIKey(r.Context(),apikey)
		if err!=nil{
			respondWithError(w,400,fmt.Sprintf("Couldnt get the user %v",err))
			return 
		}

		handler(w,r,user)
	}
	
}