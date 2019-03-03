package background

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

type Set []Background

type Background struct {
	ID          string      `json:"id"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	Color       string      `json:"color"`
	Description interface{} `json:"description"`
	Urls        struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
	} `json:"urls"`
	Links struct {
		Self             string `json:"self"`
		HTML             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	Categories             []interface{} `json:"categories"`
	Sponsored              bool          `json:"sponsored"`
	SponsoredBy            interface{}   `json:"sponsored_by"`
	SponsoredImpressionsID interface{}   `json:"sponsored_impressions_id"`
	Likes                  int           `json:"likes"`
	LikedByUser            bool          `json:"liked_by_user"`
	CurrentUserCollections []interface{} `json:"current_user_collections"`
	User                   struct {
		ID              string      `json:"id"`
		UpdatedAt       string      `json:"updated_at"`
		Username        string      `json:"username"`
		Name            string      `json:"name"`
		FirstName       string      `json:"first_name"`
		LastName        string      `json:"last_name"`
		TwitterUsername string      `json:"twitter_username"`
		PortfolioURL    string      `json:"portfolio_url"`
		Bio             string      `json:"bio"`
		Location        interface{} `json:"location"`
		Links           struct {
			Self      string `json:"self"`
			HTML      string `json:"html"`
			Photos    string `json:"photos"`
			Likes     string `json:"likes"`
			Portfolio string `json:"portfolio"`
			Following string `json:"following"`
			Followers string `json:"followers"`
		} `json:"links"`
		ProfileImage struct {
			Small  string `json:"small"`
			Medium string `json:"medium"`
			Large  string `json:"large"`
		} `json:"profile_image"`
		InstagramUsername string `json:"instagram_username"`
		TotalCollections  int    `json:"total_collections"`
		TotalLikes        int    `json:"total_likes"`
		TotalPhotos       int    `json:"total_photos"`
		AcceptedTos       bool   `json:"accepted_tos"`
	} `json:"user"`
}

func ParseFromJSON(j io.Reader) (*Set, error) {
	body, err := ioutil.ReadAll(j)
	if err != nil {
		return nil, fmt.Errorf("could not read from body: %v", err)
	}

	log.Debugf("from request, got body: %v", string(body))

	backgroundSet := make(Set, 0)
	if err := json.Unmarshal(body, &backgroundSet); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body from json: %v", err)
	}

	return &backgroundSet, nil
}
