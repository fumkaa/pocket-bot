package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) RequestToken(redirectUri string) (string, error) {
	data, err := c.doRequestReq(redirectUri)
	if err != nil {
		return "", err
	}

	// log.Print(string(data))

	var res RequestToken
	if err := json.Unmarshal(data, &res); err != nil {
		return "", fmt.Errorf("[RequestToken]can't encode body response in RequestToken: %w", err)
	}
	return res.RequestToken, nil
}

func (c *Client) doRequestReq(redirectUri string) ([]byte, error) {
	url := url.URL{
		Scheme: "https",
		Host:   Host,
		Path:   PathRequestToken,
	}
	sentData := ParametersRequestToken{
		ConsumerKey: c.consumerKey,
		RedirectUri: redirectUri,
	}

	bodyReq, err := json.Marshal(sentData)
	if err != nil {
		return nil, fmt.Errorf("[doRequest]can't encode in ParametersRequestToken: %w", err)
	}

	// log.Print(string(bodyReq))

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(bodyReq))
	if err != nil {
		return nil, fmt.Errorf("[doRequest]can't create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")
	req.Header.Set("X-Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[doRequest]can't do request: %w", err)
	}

	defer res.Body.Close()

	bodyRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[doRequest]can't read response body: %w", err)
	}
	return bodyRes, nil
}
