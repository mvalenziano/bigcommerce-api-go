package bigcommerce

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// map checkout structs
type Checkout struct {
	ID                      string        `json:"id"`
	Cart                    Cart          `json:"cart,omitempty"`
	BillingAddress          Address       `json:"billing_address,omitempty"`
	Consignments            []interface{} `json:"consignments,omitempty"`
	Taxes                   []Tax         `json:"taxes,omitempty"`
	Coupons                 []Coupon      `json:"coupons,omitempty"`
	OrderID                 string        `json:"order_id,omitempty"`
	ShippingCostTotalIncTax float64       `json:"shipping_cost_total_inc_tax,omitempty"`
	ShippingCostTotalExTax  float64       `json:"shipping_cost_total_ex_tax,omitempty"`
	HandlingCostTotalIncTax float64       `json:"handling_cost_total_inc_tax,omitempty"`
	HandlingCostTotalExTax  float64       `json:"handling_cost_total_ex_tax,omitempty"`
	TaxTotal                float64       `json:"tax_total,omitempty"`
	SubtotalIncTax          float64       `json:"subtotal_inc_tax,omitempty"`
	SubtotalExTax           float64       `json:"subtotal_ex_tax,omitempty"`
	GrandTotal              float64       `json:"grand_total,omitempty"`
	CreatedTime             time.Time     `json:"created_time,omitempty"`
	UpdatedTime             time.Time     `json:"updated_time,omitempty"`
	CustomerMessage         string        `json:"customer_message,omitempty"`
	StaffNote               string        `json:"staff_note,omitempty"`
	Fees                    []interface{} `json:"fees,omitempty"`
}

type Tax struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type CheckoutDiscountRequest struct {
	Carts struct {
		Discounts []struct {
			DiscountedAmount float64 `json:"discounted_amount"`
			Name             string  `json:"name"`
		} `json:"discounts"`
	} `json:"carts"`
}

// GetCart gets a cart by ID from BigCommerce and returns it
func (bc *Client) GetCheckout(checkoutID string) (*Checkout, error) {
	req := bc.getAPIRequest(http.MethodGet, "/v3/checkouts/"+checkoutID+"?include=consignments.available_shipping_options", nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := processBody(res)
	if err != nil {
		return nil, err
	}
	var cartResponse struct {
		Data Checkout `json:"data,omitempty"`
		Meta struct {
		} `json:"meta,omitempty"`
	}
	err = json.Unmarshal(b, &cartResponse)
	if err != nil {
		return nil, err
	}
	return &cartResponse.Data, nil
}

func (bc *Client) AddDiscountToCheckout(checkoutID string, discountAmount float64, discountName string) (*Cart, error) {
	var bodyRequestStruct CheckoutDiscountRequest
	bodyRequestStruct.Carts.Discounts[0].DiscountedAmount = discountAmount
	bodyRequestStruct.Carts.Discounts[0].Name = discountName

	var body []byte
	body, _ = json.Marshal(bodyRequestStruct)
	req := bc.getAPIRequest(http.MethodPost, "/v3/checkouts/"+checkoutID+"/discounts", bytes.NewReader(body))
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := processBody(res)
	if err != nil {
		return nil, err
	}
	var cartResponse struct {
		Data Cart `json:"data,omitempty"`
		Meta struct {
		} `json:"meta,omitempty"`
	}
	err = json.Unmarshal(b, &cartResponse)
	if err != nil {
		return nil, err
	}
	return &cartResponse.Data, nil
}
