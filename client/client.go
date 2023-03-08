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
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/astarte-platform/astarte-go/astarteservices"
	"github.com/astarte-platform/astarte-go/auth"
)

const defaultJWTExpiry = 300

type Client struct {
	baseURL            *url.URL
	appEngineURL       *url.URL
	housekeepingURL    *url.URL
	pairingURL         *url.URL
	realmManagementURL *url.URL
	userAgent          string
	httpClient         *http.Client
	token              string
	privateKey         []byte
	expiry             int
}

type Option = func(c *Client) error

// Finally, generics (actually, type constraints)
type privateKeyProvider interface {
	string | []byte
}

// The New function creates a new Astarte API client.
// You must provide at least an Astarte base URL and an auth resource (JWT or private key).
// If no other options are specified, the following is assumed:
// - standard Astarte URL hierarchy
// - standard HTTP client
// - "astarte-go" as user agent
func New(options ...Option) (*Client, error) {
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

// The WithAppEngineURL function allows to specify an
// AppEngine URL different from the standard one (e.g. http://localhost:4000).
// This is not recommendend in production.
func WithAppEngineURL(appEngineURL string) Option {
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
func WithHousekeepingURL(housekeepingURL string) Option {
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
func WithPairingURL(pairingURL string) Option {
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
func WithRealmManagementURL(realmManagementURL string) Option {
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
func WithBaseURL(baseURL string) Option {
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
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// The WithJWT function allows to specify a JWT
// token that the client will use to interact with Astarte.
func WithJWT(token string) Option {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

// The WithUserAgent function allows to specify the User Agent
// that the client will use when making http requests.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// The WithPrivateKey function allows to specify a realm private key,
// used internally to generate a valid JWT token to all Astarte APIs with 5 minutes expiry.
// The client will use that token to interact with Astarte.
// You can provide either a path (a string) to the key file, or the key itself (a []byte).
func WithPrivateKey[T privateKeyProvider](privateKey T) Option {
	return func(c *Client) error {
		switch k := any(privateKey).(type) {
		case string:
			var err error
			c.privateKey, err = os.ReadFile(k)
			return err
		case []byte:
			c.privateKey = k
			return nil
		default:
			return ErrNoPrivateKeyProvided
		}
	}
}

// The WithExpiry function allows to specify the expiry (in seconds) for the generated
// JWT token used internally for communication with all Astarte APIs.
// The expiry must be less than 5 minutes.
func WithExpiry(expirySeconds int) Option {
	return func(c *Client) error {
		if defaultJWTExpiry < expirySeconds {
			return ErrTooHighExpiry
		}

		c.expiry = expirySeconds
		return nil
	}
}

func (c *Client) GetPairingURL() (ret *url.URL) {
	ret, _ = url.Parse(c.pairingURL.String())
	return
}

func (c *Client) GetHousekeepingURL() (ret *url.URL) {
	ret, _ = url.Parse(c.housekeepingURL.String())
	return
}

func (c *Client) GetAppengineURL() (ret *url.URL) {
	ret, _ = url.Parse(c.appEngineURL.String())
	return
}

func (c *Client) GetRealmManagementURL() (ret *url.URL) {
	ret, _ = url.Parse(c.realmManagementURL.String())
	return
}

// nolint:gocognit
func validate(c *Client) error {
	if c.baseURL != nil && (c.appEngineURL != nil || c.realmManagementURL != nil || c.housekeepingURL != nil || c.pairingURL != nil) {
		return ErrConflictingUrls
	}
	if c.baseURL == nil && c.appEngineURL == nil && c.realmManagementURL == nil && c.housekeepingURL == nil && c.pairingURL == nil {
		return ErrNoUrlsProvided
	}
	if c.token != "" && c.privateKey != nil {
		return ErrBothJWTAndPrivateKey
	}
	if c.token == "" && c.privateKey == nil {
		return ErrNoAuthProvided
	}
	if c.privateKey == nil && c.expiry != 0 {
		return ErrExpiryButNoPrivateKeyProvided
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
		c.appEngineURL, _ = url.Parse(c.baseURL.String() + "/appengine")
		c.housekeepingURL, _ = url.Parse(c.baseURL.String() + "/housekeeping")
		c.pairingURL, _ = url.Parse(c.baseURL.String() + "/pairing")
		c.realmManagementURL, _ = url.Parse(c.baseURL.String() + "/realmmanagement")
	}

	if c.expiry == 0 {
		c.expiry = defaultJWTExpiry
	}

	return c
}

func (c *Client) getJWT() string {
	// Add all types
	servicesAndClaims := map[astarteservices.AstarteService][]string{
		astarteservices.AppEngine:       {},
		astarteservices.Channels:        {},
		astarteservices.Flow:            {},
		astarteservices.Housekeeping:    {},
		astarteservices.Pairing:         {},
		astarteservices.RealmManagement: {},
	}
	if c.token == "" {
		// if we're here, we can safely assume that the key was OK
		token, _ := auth.GenerateAstarteJWTFromPEMKey(c.privateKey, servicesAndClaims, int64(c.expiry))
		return token
	}
	return c.token
}
