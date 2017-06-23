package services

import (
	"net/http"
	"time"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"errors"
)

type Logo struct {
	Name string `json:"name"`
	Url  string `json:"logoURL"`
}

type LogoService struct {
	client *http.Client
}

const LOGOS_API_URL string = "https://logos-api.funkreich.de/"

func (service *LogoService) Search(query string) (*Logo, error) {
	endpoint := LOGOS_API_URL + "?q=" + url.QueryEscape(query)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "OpenSourcery/1.0")

	response, err := service.client.Get(endpoint)
	payload, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	var logos []Logo
	if err := json.Unmarshal(payload, &logos); err != nil {
		return nil, err
	}

	if len(logos) >= 1 {
		return &logos[0], nil
	}
	return nil, errors.New("No results")
}

func NewLogoService() *LogoService {
	return &LogoService{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
