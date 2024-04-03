package bigcommerce

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type TaxZone struct {
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name"`
	Enabled               bool                   `json:"enabled"`
	ShopperTargetSettings []ShopperTargetSetting `json:"shopper_target_settings"`
}

type ShopperTargetSetting struct {
	Locations      []Location `json:"locations"`
	CustomerGroups []int      `json:"customer_groups"`
}

type Location struct {
	CountryCode      string   `json:"country_code"`
	SubdivisionCodes []string `json:"subdivision_codes"`
	PostalCodes      []string `json:"postal_codes"`
}

type TaxClassRate struct {
	ClassRates []ClassRate `json:"class_rates"`
}

type ClassRate struct {
	Rate       float32 `json:"rate"`
	TaxClassID int     `json:"tax_class_id"`
}

// GetTaxZones returns tax zones based on zoneIds
func (bc *Client) GetTaxZones(zoneIds []int) (*[]TaxZone, error) {
	// convert zoneIds to string comma separated
	var fpart string
	fpart = ""
	if len(zoneIds) > 0 {
		fpart = "?id:in=" + strconv.Itoa(zoneIds[0])
		for i := 1; i < len(zoneIds); i++ {
			fpart += "," + strconv.Itoa(zoneIds[i])
		}
	}

	url := "/v3/tax/zones" + fpart

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var taxZoneResponse struct {
		Data []TaxZone `json:"data"`
	}
	err = json.Unmarshal(body, &taxZoneResponse)
	if err != nil {
		return nil, err
	}

	return &taxZoneResponse.Data, nil
}

func (bc *Client) GetTaxRates(taxZoneIds []int) (*[]TaxClassRate, error) {
	// convert zoneIds to string comma separated
	var fpart string
	fpart = ""
	if len(taxZoneIds) > 0 {
		fpart = "?tax_zone_id:in=" + strconv.Itoa(taxZoneIds[0])
		for i := 1; i < len(taxZoneIds); i++ {
			fpart += "," + strconv.Itoa(taxZoneIds[i])
		}
	}

	url := "/v3/tax/rates" + fpart

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		return nil, err
	}

	var taxRatesResponse struct {
		Data []TaxClassRate `json:"data"`
	}
	err = json.Unmarshal(body, &taxRatesResponse)
	if err != nil {
		return nil, err
	}

	return &taxRatesResponse.Data, nil
}
