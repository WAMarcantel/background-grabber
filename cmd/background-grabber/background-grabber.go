package main

import (
	"background-grabber/internal/grabber"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	config *grabber.Config
	logLevel string
)

func init() {
	config = grabber.NewConfig()

	flag.StringVar(&logLevel, "logLevel", "info", "log level - acceptable levels: panic, fatal, error, warn or warning, info, debug, trace")
	flag.StringVar(&config.AccessKey, "accessKey", "", "Access key to your Unsplash Account")
	flag.BoolVar(&config.Featured, "featured", false, "Whether or not to limit an update to featured photos from Unsplash.")
	flag.StringVar(&config.Collections, "collections", "", "A comma-separated list of collection IDs to filter on. The update will only return items in the specified collections.")
	flag.StringVar(&config.BackgroundsDirPath, "backgroundsDirPath", "", "Path to the backgrounds directory on your machine.")
	flag.StringVar(&config.Username, "username", "", "Limit selection to a single user.")
	flag.StringVar(&config.Query, "query", "", "Limit selection to photos matching a search term.")
	flag.StringVar(&config.Orientation, "orientation", "", "Filter search results by photo orientation. Valid values are landscape, portrait, and squarish.")
	flag.IntVar(&config.Count, "count", 1, "The number of photos to return. (Default: 1; max: 30)")
	flag.IntVar(&config.RefreshMinutes, "refreshMinutes", 60*24, "The number of minutes this program will wait until it refreshes your backgrounds")

	flag.Parse()
}

func main() {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil{
		panic(err)
	}
	log.SetLevel(lvl)

	if !config.GoodFlags() {
		flag.PrintDefaults()
		os.Exit(-1)
	} else {
		if err := run(); err != nil {
			panic(err)
		}
	}

}

func run() error {
	return grabber.New(config).Run()
}