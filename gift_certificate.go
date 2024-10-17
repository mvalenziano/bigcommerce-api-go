package bigcommerce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// map gift certificate structs
type GiftCertificate struct {
	ID           int       `json:"id,omitempty"`
	Code         string    `json:"code"`
	Amount       string    `json:"amount"`
	Status       string    `json:"status"`
	Balance      string    `json:"balance"`
	ToName       string    `json:"to_name,omitempty"`
	OrderID      int       `json:"order_id"`
	Template     string    `json:"template,omitempty"`
	Message      string    `json:"message"`
	ToEmail      string    `json:"to_email,omitempty"`
	FromName     string    `json:"from_name,omitempty"`
	FromEmail    string    `json:"from_email,omitempty"`
	CustomerID   int       `json:"customer_id,omitempty"`
	ExpiryDate   time.Time `json:"expiry_date"`   // could be time.Time if parsed
	PurchaseDate time.Time `json:"purchase_date"` // could be time.Time if parsed
	CurrencyCode string    `json:"currency_code"`
}

// GetCart gets a cart by ID from BigCommerce and returns it
func (bc *Client) GetGiftCertificateByCode(code string) (*GiftCertificate, error) {
	url := fmt.Sprintf("/v2/gift_certificates?code=%s", code)

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var giftCertificateResponse []GiftCertificate

	err = json.Unmarshal(b, &giftCertificateResponse)
	if err != nil {
		return nil, err
	}

	if len(giftCertificateResponse) == 0 {
		return nil, nil
	}

	return &giftCertificateResponse[0], nil
}

// create a gift certificate
func (bc *Client) CreateGiftCertificate(giftCertificate *GiftCertificate) (*GiftCertificate, error) {
	body, _ := json.Marshal(giftCertificate)

	log.Println("Creating gift certificate with body: ", string(body))

	req := bc.getAPIRequest(http.MethodPost, "/v2/gift_certificates", bytes.NewReader(body))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var giftCertificateResponse GiftCertificate
	err = json.Unmarshal(b, &giftCertificateResponse)
	if err != nil {
		return nil, err
	}

	return &giftCertificateResponse, nil
}

// update a gift certificate
func (bc *Client) UpdateGiftCertificate(giftCertificate *GiftCertificate) (*GiftCertificate, error) {
	body, _ := json.Marshal(giftCertificate)
	req := bc.getAPIRequest(http.MethodPut, fmt.Sprintf("/v2/gift_certificates/%d", giftCertificate.ID), bytes.NewReader(body))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var giftCertificateResponse GiftCertificate
	err = json.Unmarshal(b, &giftCertificateResponse)
	if err != nil {
		return nil, err
	}

	return &giftCertificateResponse, nil
}
