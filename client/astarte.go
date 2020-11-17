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
	"errors"
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

// Exported errors
var (
	ErrMalformedPayload = errors.New("received an invalid JSONAPI payload")
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

// Links is a struct that represent the links metadata returned by Astarte API.
// This metadata is used in Astarte APIs to perform pagination, allowing the
// user to simply follow the Next link, if any, to fetch the next page
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
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
// Only services added in the individualURLs map will be instantiated - the others will be nil
func NewClientWithIndividualURLs(individualURLs map[misc.AstarteService]string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	c := &Client{httpClient: httpClient, baseURL: nil, UserAgent: userAgent}

	for k, v := range individualURLs {
		// Parse URL
		parsedURL, err := url.Parse(v)
		if err != nil {
			return nil, err
		}

		switch k {
		case misc.AppEngine:
			c.AppEngine = &AppEngineService{client: c, appEngineURL: parsedURL}
		case misc.Housekeeping:
			c.Housekeeping = &HousekeepingService{client: c, housekeepingURL: parsedURL}
		case misc.Pairing:
			c.Pairing = &PairingService{client: c, pairingURL: parsedURL}
		case misc.RealmManagement:
			c.RealmManagement = &RealmManagementService{client: c, realmManagementURL: parsedURL}
		}
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

// SetTokenFromPrivateKeyFile generates a token from the supplied private key file and uses it for the session.
// The token will have complete API access and won't expire. To limit this behavior, either use
// SetTokenFromPrivateKeyFileWithTTL or SetTokenFromPrivateKeyFileWithClaims
func (c *Client) SetTokenFromPrivateKeyFile(privateKeyFile string) error {
	return c.SetTokenFromPrivateKeyFileWithTTL(privateKeyFile, 0)
}

// SetTokenFromPrivateKeyFileWithTTL generates a token from the supplied private key file and uses it for the session.
// The token will have complete API access and will expire in `ttlSeconds`. To further limit this behavior,
// use SetTokenFromPrivateKeyFileWithClaims
func (c *Client) SetTokenFromPrivateKeyFileWithTTL(privateKeyFile string, ttlSeconds int64) error {
	// Add all types
	servicesAndClaims := map[misc.AstarteService][]string{
		misc.AppEngine:       {},
		misc.Channels:        {},
		misc.Flow:            {},
		misc.Housekeeping:    {},
		misc.Pairing:         {},
		misc.RealmManagement: {},
	}
	return c.SetTokenFromPrivateKeyFileWithClaims(privateKeyFile, servicesAndClaims, ttlSeconds)
}

// SetTokenFromPrivateKeyFileWithClaims generates a token from the supplied private key file and uses it for the session.
// The token will have API access defined by the `servicesAndClaims` scope and will expire in `ttlSeconds`
func (c *Client) SetTokenFromPrivateKeyFileWithClaims(privateKeyFile string, servicesAndClaims map[misc.AstarteService][]string, ttlSeconds int64) error {
	var err error
	c.token, err = misc.GenerateAstarteJWTFromKeyFile(privateKeyFile, servicesAndClaims, ttlSeconds)
	return err
}

// SetTokenFromPrivateKey generates a token from the supplied private key and uses it for the session.
// The token will have complete API access and won't expire. To limit this behavior, either use
// SetTokenFromPrivateKeyWithTTL or SetTokenFromPrivateKeyWithClaims
func (c *Client) SetTokenFromPrivateKey(privateKey []byte) error {
	return c.SetTokenFromPrivateKeyWithTTL(privateKey, 0)
}

// SetTokenFromPrivateKeyWithTTL generates a token from the supplied private key and uses it for the session.
// The token will have complete API access and will expire in `ttlSeconds`. To further limit this behavior,
// use SetTokenFromPrivateKeyWithClaims
func (c *Client) SetTokenFromPrivateKeyWithTTL(privateKey []byte, ttlSeconds int64) error {
	// Add all types
	servicesAndClaims := map[misc.AstarteService][]string{
		misc.AppEngine:       {},
		misc.Channels:        {},
		misc.Flow:            {},
		misc.Housekeeping:    {},
		misc.Pairing:         {},
		misc.RealmManagement: {},
	}
	return c.SetTokenFromPrivateKeyWithClaims(privateKey, servicesAndClaims, ttlSeconds)
}

// SetTokenFromPrivateKeyWithClaims generates a token from the supplied private key and uses it for the session.
// The token will have API access defined by the `servicesAndClaims` scope and will expire in `ttlSeconds`
func (c *Client) SetTokenFromPrivateKeyWithClaims(privateKey []byte, servicesAndClaims map[misc.AstarteService][]string, ttlSeconds int64) error {
	var err error
	c.token, err = misc.GenerateAstarteJWTFromPEMKey(privateKey, servicesAndClaims, ttlSeconds)
	return err
}

// SetToken sets a JWT Token to be used by the client to authenticate. If you don't have a token, but rather
// you have a Private Key, you can use the SetTokenFromPrivateKey helper functions
func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) genericJSONDataAPIGET(ret interface{}, urlString string, expectedReturnCode int) error {
	return c.genericJSONDataAPIGETWithLinks(ret, nil, urlString, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIGETWithLinks(ret interface{}, retLinks *Links, urlString string, expectedReturnCode int) error {
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return c.doJSONAPIReqWithLinks(ret, retLinks, req, expectedReturnCode)
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

func (c *Client) genericJSONDataAPIPostWithResponse(ret interface{}, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteWithResponse(ret, "POST", urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPutWithResponse(ret interface{}, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteWithResponse(ret, "PUT", urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIPatchWithResponse(ret interface{}, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteWithResponseWithContentType(ret, "PATCH", urlString, dataPayload, "application/merge-patch+json", expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteNoResponse(httpVerb string, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWrite(nil, httpVerb, urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteWithResponse(ret interface{}, httpVerb string, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWrite(ret, httpVerb, urlString, dataPayload, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteNoResponseWithContentType(httpVerb string, urlString string, dataPayload interface{},
	contentType string, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteWithContentType(nil, httpVerb, urlString, dataPayload, contentType, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteWithResponseWithContentType(ret interface{}, httpVerb string, urlString string, dataPayload interface{},
	contentType string, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteWithContentType(ret, httpVerb, urlString, dataPayload, contentType, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWrite(ret interface{}, httpVerb string, urlString string, dataPayload interface{}, expectedReturnCode int) error {
	return c.genericJSONDataAPIWriteWithContentType(ret, httpVerb, urlString, dataPayload, "application/json", expectedReturnCode)
}

func (c *Client) genericJSONDataAPIWriteWithContentType(ret interface{}, httpVerb string, urlString string, dataPayload interface{},
	contentType string, expectedReturnCode int) error {
	var requestBody struct {
		Data interface{} `json:"data"`
	}
	requestBody.Data = dataPayload

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(httpVerb, urlString, b)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return c.doJSONAPIReq(ret, req, expectedReturnCode)
}

func (c *Client) genericJSONDataAPIDelete(urlString string, expectedReturnCode int) error {
	req, err := http.NewRequest("DELETE", urlString, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", c.UserAgent)

	return c.doJSONAPIReq(nil, req, expectedReturnCode)
}

func (c *Client) doJSONAPIReq(ret interface{}, req *http.Request, expectedReturnCode int) error {
	return c.doJSONAPIReqWithLinks(ret, nil, req, expectedReturnCode)
}

func (c *Client) doJSONAPIReqWithLinks(ret interface{}, retLinks *Links, req *http.Request, expectedReturnCode int) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedReturnCode {
		return errorFromJSONErrors(resp.Body)
	}

	// If we don't want the reply, discard the body and return
	if ret == nil {
		_, err := io.Copy(ioutil.Discard, resp.Body)
		return err
	}

	// Parse the payload as we should. This means we have to look for the
	// "data" enclosure for data and "links" for links.
	decoder := json.NewDecoder(resp.Body)

	foundData := false
	// We initialize it like this so it's already true if retLinks is nil
	foundLinks := retLinks == nil
	// Iterate until we decode data and links (if we need them). If we find
	// something wrong, return an error
	for t, err := decoder.Token(); err != io.EOF; t, err = decoder.Token() {
		if err != nil {
			// Errors in decoding, return
			return err
		}

		switch t {
		case "data":
			if err = decoder.Decode(ret); err != nil {
				return err
			}
			foundData = true
		case "links", retLinks != nil:
			if err = decoder.Decode(retLinks); err != nil {
				return err
			}
			foundLinks = true
		}

		if foundData && foundLinks {
			// We're done
			return nil
		}
	}

	// If we got here, we didn't find all that we needed
	return ErrMalformedPayload
}
