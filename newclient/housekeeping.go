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
	"fmt"
	"net/http"
	"net/url"
	"path"

	"moul.io/http2curl"
)

type ListRealmsRequest struct {
	req     *http.Request
	expects int
}

// ListRealms builds a request to list all realms in the cluster.
func (c *Client) ListRealms() (AstarteRequest, error) {
	callURL, _ := url.Parse(c.housekeepingURL.String())
	callURL.Path = path.Join(callURL.Path, "/v1/realms")

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return ListRealmsRequest{req: req, expects: 200}, nil
}

func (r ListRealmsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return ListRealmsResponse{res: res}, nil
}

func (r ListRealmsRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type GetRealmRequest struct {
	req     *http.Request
	expects int
}

// GetRealm builds a request to get data about a single Realm.
func (c *Client) GetRealm(realm string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.housekeepingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/realms/%s", realm))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return GetRealmRequest{req: req, expects: 200}, nil
}

func (r GetRealmRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return GetRealmResponse{res: res}, nil
}

func (r GetRealmRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type CreateRealmRequest struct {
	req     *http.Request
	expects int
}

type newRealmRequestBuilder struct {
	realmName                    string         `json:"realm_name"`
	publicKey                    string         `json:jwt_public_key_pem`
	replicationFactor            int            `json:replication_factor,omitempty`
	datacenterReplicationFactors map[string]int `json:datacenter_replication_factors,omitempty`
	replicationClass             string         `json:replication_class,omitempty`
}

type realmOption func(*newRealmRequestBuilder)

// CreateRealm builds a request to create a new Realm in the Cluster with default parameters.
// When running in production, it is advised to use a NetworkTopologyStrategy, or at least a
// replication factor > 1.
// You can create a realm with:
// c.NewRealm(client.WithRealmName("test"), client.WithRealmPublicKey("YOUR_REALM_PUBLIC_KEY"), client.WithReplicationFactor(3))
func (c *Client) CreateRealm(opts ...realmOption) (AstarteRequest, error) {
	newRealm := newRealmRequestBuilder{}
	for _, f := range opts {
		f(&newRealm)
	}

	if err := newRealm.validate(); err != nil {
		return Empty{}, err
	}

	// TODO check if setting default value is needed
	// if newRealm.datacenterReplicationFactors == nil && newRealm.replicationFactor == 0 {
	// 	newRealm.replicationFactor = 1
	// }

	callURL, _ := url.Parse(c.housekeepingURL.String())
	callURL.Path = path.Join(callURL.Path, "/v1/realms")

	// TODO check error
	reqBody, _ := makeBody(newRealm)
	req := c.makeHTTPrequest(http.MethodPost, callURL, reqBody, c.token)
	return CreateRealmRequest{req: req, expects: 201}, nil
}

func (r *newRealmRequestBuilder) validate() error {
	if r.realmName == "" {
		return ErrRealmNameNotProvided
	}
	if r.publicKey == "" {
		return ErrRealmNameNotProvided
	}
	if r.replicationFactor != 0 && r.datacenterReplicationFactors != nil {
		return ErrTooManyReplicationFactors
	}
	if r.datacenterReplicationFactors == nil && r.replicationFactor < 0 {
		return ErrNegativeReplicationFactor
	}
	return nil
}

// Sets the name for a new Realm.
func WithRealmName(name string) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.realmName = name
	}
}

// Sets the public key for a new Realm.
func WithRealmPublicKey(publicKey string) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.publicKey = publicKey
	}
}

// Sets the Replication factor for a new Realm in a single datacenter.
// Production-ready deployments usually are replicated on more datacenters,
// but if you need to use just one, set a value at least higher than 1.
func WithReplicationFactor(replicationFactor int) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.replicationFactor = replicationFactor
		req.replicationClass = fmt.Sprintf("\"SimpleStrategy\"")
	}
}

// Sets the per-datacenter Replication Factor for a new realm. This is the way to go for production deployments.
func WithDatacenterReplicationFactors(datacenterReplicationFactors map[string]int) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.datacenterReplicationFactors = datacenterReplicationFactors
		req.replicationClass = fmt.Sprintf("\"NetworkTopologyStrategy\"")
	}
}

func (r CreateRealmRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return CreateRealmResponse{res: res}, nil
}

func (r CreateRealmRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
