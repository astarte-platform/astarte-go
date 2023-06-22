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
	"fmt"
	"net/http"

	"moul.io/http2curl"
)

type ListRealmsRequest struct {
	req     *http.Request
	expects int
}

// ListRealms builds a request to list all realms in the cluster.
func (c *Client) ListRealms() (AstarteRequest, error) {
	callURL := makeURL(c.housekeepingURL, "/v1/realms")
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return ListRealmsRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r ListRealmsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListRealmsResponse{res: res}, nil
}

func (r ListRealmsRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type GetRealmRequest struct {
	req     *http.Request
	expects int
}

// GetRealm builds a request to get data about a single Realm.
func (c *Client) GetRealm(realm string) (AstarteRequest, error) {
	callURL := makeURL(c.housekeepingURL, "/v1/realms/%s", realm)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return GetRealmRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r GetRealmRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return GetRealmResponse{res: res}, nil
}

func (r GetRealmRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type CreateRealmRequest struct {
	req     *http.Request
	expects int
}

type newRealmRequestBuilder struct {
	RealmName                    string         `json:"realm_name"`
	PublicKey                    string         `json:"jwt_public_key_pem"`
	ReplicationFactor            int            `json:"replication_factor,omitempty"`
	DatacenterReplicationFactors map[string]int `json:"datacenter_replication_factors,omitempty"`
	ReplicationClass             string         `json:"replication_class,omitempty"`
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

	// TODO check if setting default replicationFactor is needed

	callURL := makeURL(c.housekeepingURL, "/v1/realms")
	reqBody, _ := makeBody(newRealm)
	req := c.makeHTTPrequest(http.MethodPost, callURL, reqBody)

	return CreateRealmRequest{req: req, expects: 201}, nil
}

func (r *newRealmRequestBuilder) validate() error {
	if r.RealmName == "" {
		return ErrRealmNameNotProvided
	}
	if r.PublicKey == "" {
		return ErrRealmPublicKeyNotProvided
	}
	if r.ReplicationFactor != 0 && r.DatacenterReplicationFactors != nil {
		return ErrTooManyReplicationFactors
	}
	if r.DatacenterReplicationFactors == nil && r.ReplicationFactor < 0 {
		return ErrNegativeReplicationFactor
	}
	return nil
}

// Sets the name for a new Realm.
// nolint:golint,revive
func WithRealmName(name string) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.RealmName = name
	}
}

// Sets the public key for a new Realm.
// nolint:golint,revive
func WithRealmPublicKey(publicKey string) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.PublicKey = publicKey
	}
}

// Sets the Replication factor for a new Realm in a single datacenter.
// Production-ready deployments usually are replicated on more datacenters,
// but if you need to use just one, set a value at least higher than 1.
// nolint:golint,revive
func WithReplicationFactor(replicationFactor int) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.ReplicationFactor = replicationFactor
		//nolint:gosimple
		req.ReplicationClass = fmt.Sprintf("\"SimpleStrategy\"")
	}
}

// Sets the per-datacenter Replication Factor for a new realm. This is the way to go for production deployments.
// nolint:golint,revive
func WithDatacenterReplicationFactors(datacenterReplicationFactors map[string]int) realmOption {
	return func(req *newRealmRequestBuilder) {
		req.DatacenterReplicationFactors = datacenterReplicationFactors
		//nolint:gosimple
		req.ReplicationClass = fmt.Sprintf("\"NetworkTopologyStrategy\"")
	}
}

// nolint:bodyclose
func (r CreateRealmRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return CreateRealmResponse{res: res}, nil
}

func (r CreateRealmRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
