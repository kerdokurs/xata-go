package xatago

import (
	"bytes"
	"encoding/json"
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
	//Message string `json:"message,omitempty"`
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

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	//if int(resp.StatusCode/100) != 2 {
	//	err = errors.New(out.(baseResponse).Message)
	//}

	return err
}

func (c *Client) query(tableName string, query *Query) ([]any, error) {
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
	return res.Records, c.doRequest(req, &res)
}

func (c *Client) create(tableName string, data any) (string, error) {
	type response struct {
		baseResponse
		ID string `json:"id"`
	}

	url := fmt.Sprintf("%s/tables/%s/data", c.DatabaseURL, tableName)

	req, err := c.buildRequest("POST", url, data)
	if err != nil {
		return "", err
	}

	res := response{}
	return res.ID, c.doRequest(req, &res)
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
	return c.doRequest(req, &res)
}
