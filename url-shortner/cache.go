package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var urlCache = cache.New(
	10*time.Minute, // default TTL
	15*time.Minute, // cleanup interval
)


func GetCachedUrl(id string)(string,bool)  {

	val,found:=urlCache.Get(id)

	if !found{
		return "",false
	}

	return val.(string),found
	
}

func setCachedUrl(id,url string){
	urlCache.Set(id,url,cache.DefaultExpiration)
}