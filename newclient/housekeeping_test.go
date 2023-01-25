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

package newclient

import (
	"testing"
)

func TestListRealms(t *testing.T) {
	c, _ := getTestContext(t)
	listRealmsCall, err := c.ListRealms()
	if err != nil {
		t.Error(err)
	}
	res, err := listRealmsCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	realms, _ := data.([]string)
	for i := 0; i < len(testRealmsList); i++ {
		if realms[i] != testRealmsList[i] {
			t.Errorf("Listed realms not matching: %s vs %s", realms[i], testRealmsList[i])
		}
	}
}

func TestGetRealm(t *testing.T) {
	c, _ := getTestContext(t)
	getRealmCall, err := c.GetRealm(testRealmName)
	if err != nil {
		t.Error(err)
	}
	res, err := getRealmCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	dat, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	details, _ := dat.(RealmDetails)
	if details.Name != testRealmName || details.JwtPublicKeyPEM != testPublicKey || details.ReplicationFactor != testReplicationFactor {
		t.Error("Received invalid realm details")
	}
}

func TestCreateRealm(t *testing.T) {
	c, _ := getTestContext(t)
	createRealmCall, err := c.CreateRealm(
		WithRealmName(testRealmName),
		WithRealmPublicKey(testPublicKey),
		WithReplicationFactor(testReplicationFactor),
	)
	if err != nil {
		t.Error(err)
	}
	res, err := createRealmCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	dat, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	details, _ := dat.(RealmDetails)
	if details.Name != testRealmName || details.JwtPublicKeyPEM != testPublicKey || details.ReplicationFactor != testReplicationFactor {
		t.Error("Failed realm creations, different realm details")
	}
}
