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
	"encoding/json"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// Parses data obtained by performing a request to list realms.
// Returns the list of realms as an array of strings.
func (r ListRealmsResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	body := string(b)
	ret := []string{}
	for _, v := range gjson.Get(body, "data").Array() {
		ret = append(ret, v.Str)
	}
	return ret, nil
}
func (r ListRealmsResponse) Raw() *http.Response {
	return r.res
}

// RealmDetails represents details of a single Realm.
type RealmDetails struct {
	Name                         string         `json:"realm_name"`
	JwtPublicKeyPEM              string         `json:"jwt_public_key_pem"`
	ReplicationClass             string         `json:"replication_class,omitempty"`
	ReplicationFactor            int            `json:"replication_factor,omitempty"`
	DatacenterReplicationFactors map[string]int `json:"datacenter_replication_factors,omitempty"`
}

// Parses data obtained by performing a request to get a realm's details.
// Returns the details as a RealmDetails struct.
func (r GetRealmResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := RealmDetails{}
	// TODO check err
	_ = json.Unmarshal(v, &ret)
	return ret, nil

}
func (r GetRealmResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to create a realm.
// Returns the realm's details as a RealmDetails struct.
func (r CreateRealmResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := RealmDetails{}
	// TODO check err
	_ = json.Unmarshal(v, &ret)
	return ret, nil
}
func (r CreateRealmResponse) Raw() *http.Response {
	return r.res
}
