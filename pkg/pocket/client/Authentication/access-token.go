package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (c *Client) AccessToken(requestToken string) (string, error) {
	data, err := c.doRequestAccess(requestToken)
	if err != nil {
		return "", fmt.Errorf("[AccessToken]can't get body response: %w", err)
	}
	log.Print(string(data))
	var res AccessToken
	if err := json.Unmarshal(data, &res); err != nil {
		return "", fmt.Errorf("[AccessToken]can't encode body response in AccessToken: %w", err)
	}
	return res.AccessToken, nil
}

func (c *Client) doRequestAccess(requestToken string) ([]byte, error) {
	url := url.URL{
		Scheme: "https",
		Host:   Host,
		Path:   PathAccessToken,
	}

	sentData := ParametersAccessToken{
		ConsumerKey:  c.consumerKey,
		RequestToken: requestToken,
	}

	bodyReq, err := json.Marshal(sentData)
	if err != nil {
		return nil, fmt.Errorf("[doRequestAccess]can't encode in ParametersRequestToken: %w", err)
	}

	// log.Print(string(bodyReq))

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(bodyReq))
	if err != nil {
		return nil, fmt.Errorf("[doRequestAccess]can't create request: %w", err)
	}

	// log.Print(req)

	req.Header.Set("Content-Type", "application/json; charset=UTF8")
	req.Header.Set("X-Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[doRequestAccess]can't do request: %w", err)
	}

	defer res.Body.Close()

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[doRequestAccess]can't read response body: %w", err)
	}
	return bodyRes, nil
}
