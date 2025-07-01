package bigcommerce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// GetAuthContext returns an AuthContext object from the BigCommerce API
// Call it with r.URL.Query() - will return BigCommerce Auth Context or error
func (bc *App) GetAuthContext(requestURLQuery url.Values) (*AuthContext, error) {

	req := AuthTokenRequest{
		ClientID:     bc.AppClientID,
		ClientSecret: bc.AppClientSecret,
		RedirectURI:  "https://" + bc.Hostname + "/auth",
		GrantType:    "authorization_code",
		Code:         requestURLQuery.Get("code"),
		Scope:        requestURLQuery.Get("scope"),
		Context:      requestURLQuery.Get("context"),
	}
	reqb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := bc.HTTPClient.Post("https://login.bigcommerce.com/oauth2/token",
		"application/json",
		bytes.NewReader(reqb),
	)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	if strings.Contains(string(bytes), "invalid_") {
		return nil, fmt.Errorf("%s", string(bytes))
	}

	var ac AuthContext
	err = json.Unmarshal(bytes, &ac)
	if err != nil {
		return nil, err
	}
	if ac.Error != "" {
		return nil, fmt.Errorf("AuthContext error: %s", ac.Error)
	}
	return &ac, nil
}
