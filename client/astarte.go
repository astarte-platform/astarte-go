// Copyright © 2019-2020 Ispirata Srl
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
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/astarte-platform/astarte-go/misc"
)

const (
	userAgent = "astarte-go"
)

// Client is the base Astarte API client. It provides access to all of Astarte's APIs.
// To use a Client, create one using the the NewClient or the NewClientWithIndividualURLs methods.
// Client will expose a set of Services each corresponding to their Astarte APIs. Please note that
// when using NewClientWithIndividualURLs, if an URL for a specific API set is not provided, the
// Service won't be available and `nil` will be returned. It is guaranteed, instead, that when
// using NewClient, all of the API Services will be allocated.
//
// Before using a Client, you must set an Authentication Token. To do so, you can invoke the
// SetToken functions, which provide a number of helper mechanisms to use Private Keys.
// You can reset the token at any time, and it will be evaluated before every API invocation.
// In most cases, you want to map an individual Client object to either Housekeeping or a
// Realm, but in some cases you might want to reset the token often (for example, this applies
// to methods such as GetMQTTv1ProtocolInformationForDevice and ObtainNewMQTTv1CertificateForDevice,
// which require a Device Credential Secret to be set as the token).
type Client struct {
	baseURL   *url.URL
	UserAgent string

	httpClient *http.Client
	token      string

	AppEngine       *AppEngineService
	Housekeeping    *HousekeepingService
	Pairing         *PairingService
	RealmManagement *RealmManagementService
}

// NewClient creates a new Astarte API client with standard URL hierarchies.
func NewClient(rawBaseURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{httpClient: httpClient, baseURL: baseURL, UserAgent: userAgent}

	// Apparently that's how you deep-copy the URLs.
	// We're ignoring errors here as the cross-parsing cannot fail.
	appEngineURL, _ := url.Parse(baseURL.String())
	appEngineURL.Path = path.Join(appEngineURL.Path, "appengine")
	c.AppEngine = &AppEngineService{client: c, appEngineURL: appEngineURL}

	housekeepingURL, _ := url.Parse(baseURL.String())
	housekeepingURL.Path = path.Join(housekeepingURL.Path, "housekeeping")
	c.Housekeeping = &HousekeepingService{client: c, housekeepingURL: housekeepingURL}

	pairingURL, _ := url.Parse(baseURL.String())
	pairingURL.Path = path.Join(pairingURL.Path, "pairing")
	c.Pairing = &PairingService{client: c, pairingURL: pairingURL}

	realmManagementURL, _ := url.Parse(baseURL.String())
	realmManagementURL.Path = path.Join(realmManagementURL.Path, "realmmanagement")
	c.RealmManagement = &RealmManagementService{client: c, realmManagementURL: realmManagementURL}

	return c, nil
}

// NewClientWithIndividualURLs creates a new Astarte API client with custom URL hierarchies.
// If an empty string is passed as one of the URLs, the corresponding Service will not be instantiated.
func NewClientWithIndividualURLs(rawAppEngineURL string, rawHousekeepingURL string, rawPairingURL string,
	rawRealmManagementURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	c := &Client{httpClient: httpClient, baseURL: nil, UserAgent: userAgent}

	if rawAppEngineURL != "" {
		appEngineURL, err := url.Parse(rawAppEngineURL)
		if err != nil {
			return nil, err
		}
		c.AppEngine = &AppEngineService{client: c, appEngineURL: appEngineURL}
	}

	if rawHousekeepingURL != "" {
		housekeepingURL, err := url.Parse(rawHousekeepingURL)
		if err != nil {
			return nil, err
		}
		c.Housekeeping = &HousekeepingService{client: c, housekeepingURL: housekeepingURL}
	}

	if rawPairingURL != "" {
		pairingURL, err := url.Parse(rawPairingURL)
		if err != nil {
			return nil, err
		}
		c.Pairing = &PairingService{client: c, pairingURL: pairingURL}
	}

	if rawRealmManagementURL != "" {
		realmManagementURL, err := url.Parse(rawRealmManagementURL)
		if err != nil {
			return nil, err
		}
		c.RealmManagement = &RealmManagementService{client: c, realmManagementURL: realmManagementURL}
	}

	return c, nil
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

// SetTokenFromPrivateKey generates a token from the supplied private key and uses it for the session.
// The token will have complete API access and won't expire. To limit this behavior, either use
// SetTokenFromPrivateKeyWithTTL or SetTokenFromPrivateKeyWithClaims
func (c *Client) SetTokenFromPrivateKey(privateKeyFile string) error {
	return c.SetTokenFromPrivateKeyWithTTL(privateKeyFile, 0)
}

// SetTokenFromPrivateKeyWithTTL generates a token from the supplied private key and uses it for the session.
// The token will have complete API access and will expire in `ttlSeconds`. To further limit this behavior,
// use SetTokenFromPrivateKeyWithClaims
func (c *Client) SetTokenFromPrivateKeyWithTTL(privateKeyFile string, ttlSeconds int64) error {
	// Add all types
	servicesAndClaims := map[misc.AstarteService][]string{
		misc.AppEngine:       []string{},
		misc.Channels:        []string{},
		misc.Housekeeping:    []string{},
		misc.Pairing:         []string{},
		misc.RealmManagement: []string{},
	}
	return c.SetTokenFromPrivateKeyWithClaims(privateKeyFile, servicesAndClaims, ttlSeconds)
}

// SetTokenFromPrivateKeyWithClaims generates a token from the supplied private key and uses it for the session.
// The token will have API access defined by the `servicesAndClaims` scope and will expire in `ttlSeconds`
func (c *Client) SetTokenFromPrivateKeyWithClaims(privateKeyFile string, servicesAndClaims map[misc.AstarteService][]string, ttlSeconds int64) error {
	var err error
	c.token, err = misc.GenerateAstarteJWTFromKeyFile(privateKeyFile, servicesAndClaims, ttlSeconds)
	return err
}

// SetToken sets a JWT Token to be used by the client to authenticate. If you don't have a token, but rather
// you have a Private Key, you can use the SetTokenFromPrivateKey helper functions
func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) genericJSONDataAPIGET(urlString string, expectedReturnCode int) (*json.Decoder, error) {
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedReturnCode {
		return nil, errorFromJSONErrors(resp.Body)
	}

	return json.NewDecoder(resp.Body), nil
}

func (c *Client) genericJSONDataAPIPost(urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteNoResponse("POST", urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPut(urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteNoResponse("PUT", urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPatch(urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteNoResponseWithContentType("PATCH", urlString, dataPayload, "application/merge-patch+json", expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPostWithResponse(urlString string, dataPayload interface{}, expectedReturnCode int) (*json.Decoder, error) {
	return c.genericJSONDataAPIWriteWithResponse("POST", urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPutWithResponse(urlString string, dataPayload interface{}, expectedReturnCode int) (*json.Decoder, error) {
	return c.genericJSONDataAPIWriteWithResponse("PUT", urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPatchWithResponse(urlString string, dataPayload interface{}, expectedReturnCode int) (*json.Decoder, error) {
	return c.genericJSONDataAPIWriteWithResponseWithContentType("PATCH", urlString, dataPayload, "application/merge-patch+json", expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteNoResponse(httpVerb string, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	decoder, err := c.genericJSONDataAPIWrite(httpVerb, urlString, dataPayload, expectedReturnCode)
	if err != nil {
		return err
	}

	// When calling this function, we're discarding the response, but there might indeed have been
	// something in the body. To avoid screwing up our client, we need ensure the response
	// is drained and the body reader is closed.
	io.Copy(ioutil.Discard, decoder.Buffered())

	return nil
}

func (c *Client) genericJSONDataAPIWriteWithResponse(httpVerb string, urlString string, dataPayload interface{}, expectedReturnCode int) (*json.Decoder, error) {
	decoder, err := c.genericJSONDataAPIWrite(httpVerb, urlString, dataPayload, expectedReturnCode)
	if err != nil {
		return nil, err
	}

	return decoder, err
}

func (c *Client) genericJSONDataAPIWriteNoResponseWithContentType(httpVerb string, urlString string, dataPayload interface{},
	contentType string, expectedReturnCode int) error {
	decoder, err := c.genericJSONDataAPIWriteWithContentType(httpVerb, urlString, dataPayload, contentType, expectedReturnCode)
	if err != nil {
		return err
	}

	// When calling this function, we're discarding the response, but there might indeed have been
	// something in the body. To avoid screwing up our client, we need ensure the response
	// is drained and the body reader is closed.
	io.Copy(ioutil.Discard, decoder.Buffered())

	return nil
}

func (c *Client) genericJSONDataAPIWriteWithResponseWithContentType(httpVerb string, urlString string, dataPayload interface{},
	contentType string, expectedReturnCode int) (*json.Decoder, error) {
	decoder, err := c.genericJSONDataAPIWriteWithContentType(httpVerb, urlString, dataPayload, contentType, expectedReturnCode)
	if err != nil {
		return nil, err
	}

	return decoder, err
}

func (c *Client) genericJSONDataAPIWrite(httpVerb string, urlString string, dataPayload interface{}, expectedReturnCode int) (*json.Decoder, error) {
	return c.genericJSONDataAPIWriteWithContentType(httpVerb, urlString, dataPayload, "application/json", expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteWithContentType(httpVerb string, urlString string, dataPayload interface{},
	contentType string, expectedReturnCode int) (*json.Decoder, error) {
	var requestBody struct {
		Data interface{} `json:"data"`
	}
	requestBody.Data = dataPayload

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(httpVerb, urlString, b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedReturnCode {
		return nil, errorFromJSONErrors(resp.Body)
	}

	return json.NewDecoder(resp.Body), nil
}

func (c *Client) genericJSONDataAPIDelete(urlString string, expectedReturnCode int) error {
	req, err := http.NewRequest("DELETE", urlString, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != expectedReturnCode {
		return errorFromJSONErrors(resp.Body)
	}

	// When calling this function, we're discarding the response, but there might indeed have been
	// something in the body. To avoid screwing up our client, we need ensure the response
	// is drained and the body reader is closed.
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return nil
}
