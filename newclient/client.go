// Copyright © 2023 SECO Mind srl
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

package newclient

import (
	"net/http"
	"net/url"
	"time"

	"github.com/astarte-platform/astarte-go/misc"
)

type Client struct {
	baseURL            *url.URL
	appEngineURL       *url.URL
	housekeepingURL    *url.URL
	pairingURL         *url.URL
	realmManagementURL *url.URL
	userAgent          string

	httpClient *http.Client
	token      string
}

type clientOption = func(c *Client) error

// Finally, generics (actually, type constraints)
type privateKeyProvider interface {
	string | []byte
}

// The New function creates a new Astarte API client.
// If no options are specified, the following is assumed:
// - standard Astarte URL hierarchy
// - standard HTTP client
// - no JWT token (no call will be authorized)
// - "astarte-go" as user agent
// A production-ready client may be created using e.g.:
// `client.New(client.WithBaseUrl("api.your-astarte.org"), client.WithToken("YOUR_JWT_TOKEN"))``
func New(options ...clientOption) (*Client, error) {
	// We start with a client with bare zero-valued fields
	c := &Client{}

	// Then we modify it according to user-provided options...
	for _, f := range options {
		err := f(c)
		if err != nil {
			return c, err
		}
	}

	// ... and check if the result is valid
	if err := validate(c); err != nil {
		return c, err
	}

	// Finally, we add just a sprinkle of defaults and a new Client is born!
	return setDefaults(c), nil
}

// The WithAppengineURL function allows to specify an
// AppEngine URL different from the standard one (e.g. http://localhost:4000).
// This is not recommendend in production.
func WithAppengineURL(appEngineURL string) clientOption {
	return func(c *Client) error {
		appengine, err := url.Parse(appEngineURL)
		if err != nil {
			return err
		}
		c.appEngineURL = appengine
		return nil
	}
}

// The WithHousekeepingURL function allows to specify an
// Housekeeping URL different from the standard one (e.g. http://localhost:4001).
// This is not recommendend in production.
func WithHousekeepingURL(housekeepingURL string) clientOption {
	return func(c *Client) error {
		housekeeping, err := url.Parse(housekeepingURL)
		if err != nil {
			return err
		}
		c.housekeepingURL = housekeeping
		return nil
	}
}

// The WithPairingURL function allows to specify an
// Pairing URL different from the standard one (e.g. http://localhost:4002).
// This is not recommendend in production.
func WithPairingURL(pairingURL string) clientOption {
	return func(c *Client) error {
		// check that it's a valid URL
		pairing, err := url.Parse(pairingURL)
		if err != nil {
			return err
		}
		c.pairingURL = pairing
		return nil
	}
}

// The WithRealmManagementURL function allows to specify an
// RealmManagement URL different from the standard one (e.g. http://localhost:4003).
// This is not recommendend in production.
func WithRealmManagementURL(realmManagementURL string) clientOption {
	return func(c *Client) error {
		realmManagement, err := url.Parse(realmManagementURL)
		if err != nil {
			return err
		}
		c.realmManagementURL = realmManagement
		return nil
	}
}

// The WithBaseURL function allows to specify the Astarte
// base URL (e.g. api.your-astarte.org)
func WithBaseURL(baseURL string) clientOption {
	return func(c *Client) error {
		base, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.baseURL = base

		return nil
	}
}

// The WithHTTPClient function allows to specify an httpClient
// with custom options, e.g. different timeout, or skipTLSVerify
func WithHTTPClient(httpClient *http.Client) clientOption {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// The WithToken function allows to specify a JWT
// token that the client will use to interact with Astarte.
func WithToken(token string) clientOption {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

// The WithUserAgent function allows to specify the User Agent
// that the client will use when making http requests.
func WithUserAgent(userAgent string) clientOption {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// The WithPrivateKey function allows to specify a realm private key,
// used internally to generate a valid JWT token to all Astarte APIs with no expiry.
// The client will use that token to interact with Astarte.
func WithPrivateKey[T privateKeyProvider](privateKey T) clientOption {
	return WithPrivateKeyWithTTL(privateKey, 0)
}

// The WithPrivateKey function allows to specify a realm private key,
// used internally to generate a valid JWT token to all Astarte APIs with a specified expiry (in seconds).
// The client will use that token to interact with Astarte.
func WithPrivateKeyWithTTL[T privateKeyProvider](privateKey T, ttlSeconds int64) clientOption {
	// Add all types
	servicesAndClaims := map[misc.AstarteService][]string{
		misc.AppEngine:       {},
		misc.Channels:        {},
		misc.Flow:            {},
		misc.Housekeeping:    {},
		misc.Pairing:         {},
		misc.RealmManagement: {},
	}
	return WithPrivateKeyWithClaimsWithTTL(privateKey, servicesAndClaims, 0)
}

// The WithPrivateKey function allows to specify a realm private key,
// used internally to generate a valid JWT token with a given set of Astarte claims and
// a specified expiry (in seconds).
// The client will use that token to interact with Astarte.
func WithPrivateKeyWithClaimsWithTTL[T privateKeyProvider](privateKey T, claims map[misc.AstarteService][]string, ttlSeconds int64) clientOption {
	return func(c *Client) error {
		// Golang I hate you so much
		switch k := any(privateKey).(type) {
		case string:
			var err error
			c.token, err = misc.GenerateAstarteJWTFromKeyFile(k, claims, ttlSeconds)
			return err
		case []byte:
			var err error
			c.token, err = misc.GenerateAstarteJWTFromPEMKey(k, claims, ttlSeconds)
			return err
		default:
			return ErrNoPrivateKeyProvided
		}
	}
}

func validate(c *Client) error {
	if c.baseURL != nil && (c.appEngineURL != nil || c.realmManagementURL != nil || c.housekeepingURL != nil || c.pairingURL != nil) {
		return ErrConflictingUrls
	}
	if c.baseURL == nil && c.appEngineURL == nil && c.realmManagementURL == nil && c.housekeepingURL == nil && c.pairingURL == nil {
		return ErrNoUrlsProvided
	}
	return nil
}

func setDefaults(c *Client) *Client {
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: time.Second * 30,
		}

	}
	if c.userAgent == "" {
		c.userAgent = "astarte-go"
	}

	if c.baseURL != nil {
		c.appEngineURL, _ = url.Parse(c.baseURL.String()+"/appengine")
		c.housekeepingURL, _ = url.Parse(c.baseURL.String()+"/housekeeping")
		c.pairingURL, _ = url.Parse(c.baseURL.String()+"/pairing")
		c.realmManagementURL, _ = url.Parse(c.baseURL.String()+"/realmmanagement")
	}

	return c
}
