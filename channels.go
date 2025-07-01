package bigcommerce

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Channel struct {
	IconURL          string    `json:"icon_url"`
	IsListableFromUI bool      `json:"is_listable_from_ui"`
	IsVisible        bool      `json:"is_visible"`
	DateCreated      time.Time `json:"date_created"`
	ExternalID       string    `json:"external_id"`
	Type             string    `json:"type"`
	Platform         string    `json:"platform"`
	IsEnabled        bool      `json:"is_enabled"`
	DateModified     time.Time `json:"date_modified"`
	Name             string    `json:"name"`
	ID               int       `json:"id"`
	Status           string    `json:"status"`
}

func (bc *Client) GetAllChannels() ([]Channel, error) {
	cs := []Channel{}
	var csp []Channel
	page := 1
	more := true
	var err error
	retries := 0
	for more {
		csp, more, err = bc.GetChannels(page)
		if err != nil {
			retries++
			if retries > bc.MaxRetries {
				log.Println("Max retries reached")
				return cs, err
			}
			break
		}
		cs = append(cs, csp...)
		page++
	}
	return cs, err
}

func (bc *Client) GetChannels(page int) ([]Channel, bool, error) {
	url := "/v3/channels?page=" + strconv.Itoa(page)

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil, false, ErrNoContent
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, false, err
	}
	var pp struct {
		Status int       `json:"status"`
		Title  string    `json:"title"`
		Data   []Channel `json:"data"`
		Meta   struct {
			Pagination Pagination `json:"pagination"`
		} `json:"meta"`
	}
	err = json.Unmarshal(body, &pp)
	if err != nil {
		return nil, false, err
	}
	if pp.Status != 0 {
		return nil, false, errors.New(pp.Title)
	}
	return pp.Data, pp.Meta.Pagination.CurrentPage < pp.Meta.Pagination.TotalPages, nil
}
