package newclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AstarteRequest interface {
	// Run executes an astarteRequest that was built using functions from the client package.
	// To retrive the result, see the data.Parse function
	Run(c *Client) (AstarteResponse, error)
	// ToCurl returns the curl command equivalent to the provided astarteRequest.
	// This does not execute neither the request nor the command.
	ToCurl(c *Client) string
}

// This empty struct is there just for errors, method implementations are bogus
type empty struct{}

func (r empty) Run(c *Client) (AstarteResponse, error) { return Empty{}, nil }
func (r empty) ToCurl(c *Client) string                     { return "" }

func (c *Client) makeHTTPrequest(method string, url *url.URL, payload io.Reader, token string) *http.Request {
	return c.makeHTTPrequestWithContentType(method, url, payload, token, "application/json")
}

func (c *Client) makeHTTPrequestWithContentType(method string, url *url.URL, payload io.Reader, token string, contentType string) *http.Request {
	// TODO check err
	req, _ := http.NewRequest(method, url.String(), payload)
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	return req
}

type astarteRequestBody struct {
	Data any `json:"data"`
}

func makeBody(payload any) (io.Reader, error) {
	data := astarteRequestBody{Data: payload}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return b, err
	}
	return b, nil
}

func errorFromJSONErrors(responseBody io.Reader) error {
	var errorBody struct {
		Errors map[string]interface{} `json:"errors"`
	}

	err := json.NewDecoder(responseBody).Decode(&errorBody)
	if err != nil {
		return err
	}

	errJSON, _ := json.MarshalIndent(&errorBody, "", "  ")
	return fmt.Errorf("%s", errJSON)
}
