package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func PrintRed(s... interface{}){
	for _, v := range s{
		fmt.Printf("\033[91m%v\033[0m", v)
	}
}

func Connected() (ok bool) {
	_, err := http.Get("https://api.unsplash.com/")
	if err != nil {
		log.Debugf("could not connect to https://api.unsplash.com/")
		return false
	}
	log.Debugf("connected successfully to https://api.unsplash.com/")
	return true
}