package xatago

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the main Xata Client struct
type Client struct {
	httpClient  *http.Client
	accesToken  string
	DatabaseURL string
}

type baseResponse struct {
	ID      string `json:"id"`
	Message string `json:"message,omitempty"`
}

// NewClient initializes a new Xata Client
func NewClient(accessToken string, databaseURL string) *Client {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	client := &Client{
		httpClient:  httpClient,
		accesToken:  accessToken,
		DatabaseURL: databaseURL,
	}

	return client
}

func (c *Client) buildRequest(method, url string, body any) (*http.Request, error) {
	var bd io.Reader
	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bd = bytes.NewBuffer(bodyData)
	}
	req, err := http.NewRequest(method, url, bd)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accesToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "xata-go-client")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: This is a temporary fix
	if req.Method == "DELETE" && int(resp.StatusCode/100) == 2 {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	return err
}

func (c *Client) query(tableName string, query *apiQuery) ([]any, error) {
	type response struct {
		baseResponse
		//Meta    any   `json:"meta,omitempty"`
		Records []any `json:"records,omitempty"`
	}

	url := fmt.Sprintf("%s/tables/%s/query", c.DatabaseURL, tableName)

	req, err := c.buildRequest("POST", url, query)
	if err != nil {
		return nil, err
	}

	res := response{}
	if err = c.doRequest(req, &res); err != nil {
		return nil, err
	}

	if res.Message != "" {
		return nil, errors.New(res.Message)
	}

	return res.Records, nil
}

func (c *Client) create(tableName string, data, out any) error {
	url := fmt.Sprintf("%s/tables/%s/data?columns=*", c.DatabaseURL, tableName)

	var err error
	if data, err = preprocessForCreate(&data); err != nil {
		return err
	}

	req, err := c.buildRequest("POST", url, data)
	if err != nil {
		return err
	}

	resMap := make(map[string]any)
	if err = c.doRequest(req, &resMap); err != nil {
		return err
	}

	if message, exists := resMap["message"]; exists && message != "" {
		return errors.New(message.(string))
	}

	delete(resMap, "message")
	return MapToStruct[map[string]any, any](resMap, &out)
}

func (c *Client) update(tableName, id string, data, out any) error {
	url := fmt.Sprintf("%s/tables/%s/data/%s?columns=*", c.DatabaseURL, tableName, id)

	dataMap, err := preprocessForCreate(&data)
	if err != nil {
		return err
	}

	delete(dataMap, "id")

	req, err := c.buildRequest("PATCH", url, dataMap)
	if err != nil {
		return err
	}

	resMap := make(map[string]any)
	if err = c.doRequest(req, &resMap); err != nil {
		return err
	}

	if message, exists := resMap["message"]; exists && message != "" {
		return errors.New(message.(string))
	}

	delete(resMap, "message")
	return MapToStruct[map[string]any, any](resMap, &out)
}

func (c *Client) delete(tableName string, id string) error {
	url := fmt.Sprintf("%s/tables/%s/data/%s", c.DatabaseURL, tableName, id)

	type response struct {
		baseResponse
	}

	req, err := c.buildRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	res := response{}
	if err = c.doRequest(req, &res); err != nil {
		return err
	}

	if res.Message != "" {
		return errors.New(res.Message)
	}

	return nil
}

func (c *Client) GetSchema() (*Schema, error) {
	type response struct {
		Schema *Schema `json:"schema"`
	}
	req, err := c.buildRequest("GET", c.DatabaseURL, nil)
	if err != nil {
		return nil, err
	}

	res := response{}
	if err = c.doRequest(req, &res); err != nil {
		return nil, err
	}

	return res.Schema, nil
}
