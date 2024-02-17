package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) Add(url string, accessToken string) error {
	_, err := c.doRequestAdd(url, accessToken)
	if err != nil {
		return fmt.Errorf("[]can't get body response: %w", err)
	}

	return nil
}

func (c *Client) doRequestAdd(u string, accessToken string) ([]byte, error) {
	url := url.URL{
		Scheme: "https",
		Host:   Host,
		Path:   PathMethodAdd,
	}

	sentData := ParametersAdd{
		Url:         u,
		ConsumerKey: c.consumerKey,
		AccessToken: accessToken,
	}

	data, err := json.Marshal(sentData)
	if err != nil {
		return nil, fmt.Errorf("[doRequestAdd]can't encode in ParametersAdd: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("[doRequestAdd]can't create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")
	req.Header.Set("X-Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[doRequestAdd]can't do request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[doRequestAdd]can't read response: %w", err)
	}
	return body, nil
}
