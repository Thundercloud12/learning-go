package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIey(headers http.Header) (string,error){
	val:=headers.Get("Authorization")
	if val==""{
		return "",errors.New("No authoriation info nfound")
	}

	vals:=strings.Split(val," ")
	if len(vals) !=2{
		return "", errors.New("Matlformed")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("Mslformed first part")
	}

	return vals[1],nil
}