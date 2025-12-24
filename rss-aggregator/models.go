package main

import (
	"time"

	"github.com/Thundercloud12/go-projects/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey	  string	`json: "ApiKey"`
}

type Feed struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL		  string 	`json:"url"`
	User_id   uuid.UUID 	`json:"user_id"`
	
}

type FeedFollow struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

func databaseFeedtoFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		URL: 	   dbFeed.Url,
		User_id:   dbFeed.UserID,	
		
	}
}

func databaseFeedstoFeeds(dbFeeds []database.Feed) []Feed {
	feeds:= []Feed{}
	for _,dbFeed:=range dbFeeds{
		feeds=append(feeds, databaseFeedtoFeed(dbFeed))
	}

	return feeds
	
}

func databaseFeedFollowToFeedFollow(dbFeedFollow database.FeedsFollower) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

func databaseFeedUserFollortoUserFollows(dbFeedFollows []database.FeedsFollower) []FeedFollow {
	feedFollows:= []FeedFollow{}
	for _,dbFeedFollow:=range dbFeedFollows{
		feedFollows=append(feedFollows, databaseFeedFollowToFeedFollow(dbFeedFollow))
	}

	return feedFollows
	
}