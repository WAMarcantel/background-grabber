package grabber

import (
	"background-grabber/util"
	"net/url"
	"strconv"
)

type Config struct {
	AccessKey          string
	Orientation        string
	Featured           bool
	Username           string
	Query              string
	Count              int
	Collections        string
	RefreshMinutes     int
	BackgroundsDirPath string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GoodFlags() bool {

	if c.AccessKey == "" {
		util.PrintRed("must provide accessKey\n\n")
		return false
	}
	if c.Count > 30 || c.Count < 1 {
		util.PrintRed("count must be set between 1 and 30\n\n")
		return false
	}
	if !util.StringInSlice(c.Orientation, []string{"landscape", "portrait", "squarish"}) {
		util.PrintRed("orientation must be landscape, portrait, or squarish\n\n")
		return false
	}

	if c.RefreshMinutes < 1 {
		util.PrintRed("refreshMinutes must be set to a number greater than 0\n\n")
		return false
	}

	return true
}

func (c *Config) addQueryParams(u *url.URL) *url.URL {

	q := url.Values{}

	q.Add("client_id", c.AccessKey)
	q.Add("orientation", c.Orientation)
	q.Add("count", strconv.Itoa(c.Count))
	q.Add("featured", strconv.FormatBool(c.Featured))

	if c.Username != "" {
		q.Add("username", c.Username)
	}
	if c.Collections != "" {
		q.Add("collections", c.Collections)
	}
	if c.Query != "" {
		q.Add("query", c.Query)
	}

	u.RawQuery = q.Encode()
	u.Query()
	return u
}
