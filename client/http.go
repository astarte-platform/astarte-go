// Copyright © 2023 SECO Mind Srl
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type AstarteRequest interface {
	// Run executes an astarteRequest that was built using functions from this package.
	// To retrieve the result, see the Parse function.
	Run(c *Client) (AstarteResponse, error)
	// ToCurl returns the curl command equivalent to the provided astarteRequest.
	// This does not execute neither the request nor the command.
	ToCurl(_ *Client) string
}

// The Empty struct represent errors, method implementations are bogus
type Empty struct{}

func (r Empty) Run(_ *Client) (AstarteResponse, error) { return Empty{}, nil }
func (r Empty) ToCurl(_ *Client) string                { return "" }

func (c *Client) makeHTTPrequest(method string, url *url.URL, payload io.Reader) *http.Request {
	return c.makeHTTPrequestWithContentType(method, url, payload, "application/json")
}

func (c *Client) makeHTTPrequestWithContentType(method string, url *url.URL, payload io.Reader, contentType string) *http.Request {
	// TODO check err
	req, _ := http.NewRequest(method, url.String(), payload)
	req.Header.Add("Authorization", "Bearer "+c.getJWT())
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

func makeURL(base *url.URL, pathFormat string, args ...interface{}) *url.URL {
	callURL, _ := url.Parse(base.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf(pathFormat, args...))
	return callURL
}

// setupURLQuery setups URL query parameters
func setupURLQuery(u *url.URL, queries map[string]string) *url.URL {
	q := u.Query()
	for k, v := range queries {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u
}
