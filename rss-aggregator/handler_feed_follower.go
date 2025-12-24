package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/Thundercloud12/go-projects/internal/auth"
	"github.com/Thundercloud12/go-projects/internal/database"
	"github.com/google/uuid"
)


func (apiConfig  *apiConfig)handlersubscribeFeeds(w http.ResponseWriter, r *http.Request,user database.User){
	type parameters struct {
        FeedID uuid.UUID `json:"feed_id"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }
	subscribe,err:=apiConfig.DB.CreateFeedFollower(r.Context(), database.CreateFeedFollowerParams{
		ID:	uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error fetching feeds",err))
		return
	}

	respondWithJSON(w,200,(databaseFeedFollowToFeedFollow(subscribe)))

}

func (apiConfig  *apiConfig)handlergetFeedFollower(w http.ResponseWriter, r *http.Request,user database.User){
	feeds,err:= apiConfig.DB.GetFeedFollower(r.Context(),user.ID)

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error fetching feed follower",err))
		return
	}

	respondWithJSON(w,200,databaseFeedUserFollortoUserFollows(feeds))


}

func (apiConfig *apiConfig)handledeletefollower(w http.ResponseWriter, r *http.Request,user database.User){
	type parameters struct {
		ID uuid.UUID `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }
	 _, err = apiConfig.DB.ExistFeedFollower(r.Context(), database.ExistFeedFollowerParams{
        UserID: user.ID,
        ID:     params.ID,
    })
    if err != nil {
        if err == sql.ErrNoRows {

            respondWithError(w, 404, "Follower doesn't exist")
            return
        }
        respondWithError(w, 400, fmt.Sprintf("Error checking follower: %v", err))
        return
    }

	err1:=apiConfig.DB.DeleteFeedFollower(r.Context(), database.DeleteFeedFollowerParams{
		UserID: user.ID,
		ID: params.ID,
	})
	if err1 != nil {
        respondWithError(w, 400, fmt.Sprintf("Error deleting follower: %v", err))
        return
    }


	respondWithJSON(w, 200, map[string]string{"message": "Successfully unfollowed the feed"})
}