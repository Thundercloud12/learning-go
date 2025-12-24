package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Thundercloud12/go-projects/internal/database"
	"github.com/google/uuid"
)


func startScraping(db *database.Queries, concurrency int, timebetweenRequests time.Duration){

	log.Printf("Collecting feeds every %s on %v goroutines...", timebetweenRequests, concurrency)

	ticker:=time.NewTicker(timebetweenRequests)

	for ;;<-ticker.C{
		feeds,err:=db.GetNextFeedtoFetch(context.Background(),int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}

		log.Printf("Found %v feeds to fetch!", len(feeds))

		wg:= &sync.WaitGroup{}
		for _,feed:=range feeds{
			wg.Add(1)
			go scrapeFeed(db,wg,feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries,wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done()
	_,err:=db.MarkFeedAsFetched(context.Background(),feed.ID)
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}

	feedData,err:=FetchRSS(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _,item:=range feedData.Channel.Items{
		description:=sql.NullString{}
		if item.Description!= ""{
				description.String = item.Description
				description.Valid = true
		}

		pubAt,err:=time.Parse(time.RFC1123,item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse publication date %s for post %s: %v", item.PubDate, item.Title, err)
			continue
		}

		_, err1 := db.CreatePost(context.Background(),database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description.String,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err1 != nil {
			if strings.Contains(err.Error(),"duplicate key"){
				continue
			}

			log.Printf("Couldn't create post %s for feed %s: %v", item.Title, feed.Name, err)
			continue
		}

	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Items))
}