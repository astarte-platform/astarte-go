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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/astarte-platform/astarte-go/interfaces"
)

var (
	testRealmName         = "test"
	testTokenValue        = "ah yes, the token"
	testDeviceID          = "fhd0WHcgSjWeVqPGKZv_KA"
	testDeviceIDs         = []string{testDeviceID, "t1J1uQSBQRi_1F3zIrjyYw", "V_pY-ZrLQzWz4iGjGu-NuQ"}
	testBrokerUrl         = "mqtt://ah.yes.the.broker"
	testClientCrt         = "ah yes, the certificate"
	testCredentialsSecret = "ah yes, the credentials secret"
	testPublicKey         = "ah yes, the public key"
	testReplicationFactor = 3
	testRealmsList        = []string{testRealmName, "ah yes, another realm"}
	testRealmDetails      = map[string]interface{}{"realm_name": testRealmName, "jwt_public_key_pem": testPublicKey, "replication_factor": testReplicationFactor}
	testInterfacesList    = []string{"ah.yes.an.Interface", "ah.yes.another.Interface"}
	testInterfaceName     = "ah.yes.an.Interface"
	testInterfaceMajor    = 1
	testInterfaceMajors   = []int{testInterfaceMajor, 2}
	testInterfaceMinor    = 1
	testInterfaceMinors   = []int{testInterfaceMinor, 0}
	testInterface         = `{
		"interface_name": "ah.yes.an.Interface",
		"version_major": 1,
		"version_minor": 1,
		"type": "datastream",
		"ownership": "device",
		"description": "Generic sensors sampled data.",
		"doc": "Values allows generic sensors to stream samples. It is usually used in combination with AvailableSensors, which makes API client aware of what sensors and what unit of measure they are reporting. sensor_id represents an unique identifier for an individual sensor, and should match sensor_id in AvailableSensors when used in combination.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/value",
				"type": "double",
				"explicit_timestamp": true,
				"description": "Sampled real value.",
				"doc": "Datastream of sampled real values."
			}
		]
	}`
	testTriggerName  = "ah_yes_a_trigger"
	testTriggersList = []string{testTriggerName, "ah_yes_another_trigger"}
	testTrigger      = `{
		"name": "ah_yes_a_trigger",
		"action": {
			"http_post_url": "http://example.com/my_post_url"
		},
		"simple_triggers": [
			{
			"type": "device_trigger",
			"on": "device_connected",
			"device_id": "glO6LullTKmwxebForU-eg"
			}
		]
	}`
	testDevicesLinks                     = map[string]string{"self": fmt.Sprintf("/v1/%s/devices", testRealmName)}
	testServerOwnedInterfaceName         = "ah.yes.a.server.owned.Interface"
	testServerOwnedPropertyInterfaceName = "ah.yes.a.server.owned.property.Interface"
	testIndividualDatastreamSnapshot     = `
	{
		"anotherTest":{
		  "value":{
			 "reception_timestamp":"2023-01-26T15:21:38.986Z",
			 "timestamp":"2023-01-26T15:21:38.985Z",
			 "value":0.29031942518908505
		  }
		},
		"yetAnotherTest":{
		  "value":{
			 "reception_timestamp":"2023-01-26T15:23:18.485Z",
			 "timestamp":"2023-01-26T15:23:18.485Z",
			 "value":0.41505074846327805
		  }
		}
	 }
	`
	testGroupName    = "ah yes, a group"
	testGroupLinks   = map[string]string{"self": fmt.Sprintf("/v1/%s/groups/%s/devices", testRealmName, url.PathEscape(testGroupName))}
	testPolicyName   = "ah_yes_a_policy"
	testPoliciesList = []string{testPolicyName, "ah_yes_another_policy"}
	testPolicy       = `{
		"name" : "ah_yes_a_policy",
		"maximum_capacity" : 100,
		"error_handlers" : [
			{
				"on" : "any_error",
				"strategy" : "discard"
			}
		]
	  }`
)

func astarteAPIMock(w http.ResponseWriter, req *http.Request) {
	authorization := req.Header.Get("Authorization")
	if len(authorization) <= 0 {
		http.Error(w, "No token supplied", http.StatusUnauthorized)
		return
	} else if authorization != "Bearer "+testTokenValue {
		http.Error(w, "Wrong token supplied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Process request
	var reply map[string]interface{}
	switch {
	// register device
	case req.URL.Path == fmt.Sprintf("/pairing/v1/%s/agent/devices", testRealmName):
		credentialsSecret := map[string]string{"credentials_secret": testCredentialsSecret}
		reply = map[string]interface{}{"data": credentialsSecret}
		w.WriteHeader(http.StatusCreated)
	// unregister device
	case req.URL.Path == fmt.Sprintf("/pairing/v1/%s/agent/devices/%s", testRealmName, testDeviceID):
		reply = map[string]interface{}{"data": ""}
		w.WriteHeader(http.StatusNoContent)
	// get credentials
	case req.URL.Path == fmt.Sprintf("/pairing/v1/%s/devices/%s/protocols/astarte_mqtt_v1/credentials", testRealmName, testDeviceID):
		clientCrt := map[string]string{"client_crt": testClientCrt}
		reply = map[string]interface{}{"data": clientCrt}
		w.WriteHeader(http.StatusCreated)
	// get info
	case req.URL.Path == fmt.Sprintf("/pairing/v1/%s/devices/%s", testRealmName, testDeviceID):
		brokerUrl := map[string]string{"broker_url": testBrokerUrl}
		reply = map[string]interface{}{"data": brokerUrl}
	case req.URL.Path == fmt.Sprintf("/housekeeping/v1/realms"):
		if req.Method == http.MethodGet {
			// list realms
			reply = map[string]interface{}{"data": testRealmsList}
		} else if req.Method == http.MethodPost {
			// new realm
			reply = map[string]interface{}{"data": testRealmDetails}
			w.WriteHeader(http.StatusCreated)
		}
	// realm details
	case req.URL.Path == fmt.Sprintf("/housekeeping/v1/realms/%s", testRealmName):
		reply = map[string]interface{}{"data": testRealmDetails}
	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/interfaces", testRealmName):
		if req.Method == http.MethodGet {
			// interface list
			reply = map[string]interface{}{"data": testInterfacesList}
		} else if req.Method == http.MethodPost {
			// install interface
			iface, _ := interfaces.ParseInterface([]byte(testInterface))
			reply = map[string]interface{}{"data": iface}
			w.WriteHeader(http.StatusCreated)
		}
	// interface major list
	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/interfaces/%s", testRealmName, testInterfaceName):
		reply = map[string]interface{}{"data": testInterfaceMajors}

	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/interfaces/%s/%v", testRealmName, testInterfaceName, testInterfaceMajor):
		if req.Method == http.MethodGet {
			// get interface
			iface, _ := interfaces.ParseInterface([]byte(testInterface))
			reply = map[string]interface{}{"data": iface}
		} else if req.Method == http.MethodDelete {
			// delete interface
			reply = map[string]interface{}{"data": ""}
			w.WriteHeader(http.StatusNoContent)
		} else if req.Method == http.MethodPut {
			// update interface
			reply = map[string]interface{}{"data": ""}
			w.WriteHeader(http.StatusNoContent)
		}
	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/triggers", testRealmName):
		if req.Method == http.MethodGet {
			// trigger list
			reply = map[string]interface{}{"data": testTriggersList}
		} else if req.Method == http.MethodPost {
			// install trigger
			trigger := map[string]any{}
			_ = json.Unmarshal([]byte(testTrigger), &trigger)
			reply = map[string]interface{}{"data": trigger}
			w.WriteHeader(http.StatusCreated)
		}
	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/triggers/%s", testRealmName, testTriggerName):
		if req.Method == http.MethodGet {
			// get trigger
			trigger := map[string]any{}
			_ = json.Unmarshal([]byte(testTrigger), &trigger)
			reply = map[string]interface{}{"data": trigger}
		} else if req.Method == http.MethodDelete {
			// delete trigger
			reply = map[string]interface{}{"data": ""}
			w.WriteHeader(http.StatusNoContent)
		}
	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/policies", testRealmName):
		if req.Method == http.MethodGet {
			// policy list
			reply = map[string]interface{}{"data": testPoliciesList}
		} else if req.Method == http.MethodPost {
			// install policy
			policy := map[string]any{}
			_ = json.Unmarshal([]byte(testPolicy), &policy)
			reply = map[string]interface{}{"data": policy}
			w.WriteHeader(http.StatusCreated)
		}
	case req.URL.Path == fmt.Sprintf("/realmmanagement/v1/%s/policies/%s", testRealmName, testPolicyName):
		if req.Method == http.MethodGet {
			// get policy
			policy := map[string]any{}
			_ = json.Unmarshal([]byte(testPolicy), &policy)
			reply = map[string]interface{}{"data": policy}
		} else if req.Method == http.MethodDelete {
			// delete policy
			reply = map[string]interface{}{"data": ""}
			w.WriteHeader(http.StatusNoContent)
		}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/devices", testRealmName):
		reply = map[string]interface{}{"data": testDeviceIDs, "links": testDevicesLinks}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/%s/interfaces/%s", testRealmName, testDeviceID, testInterface):
		// snapshot
		data := map[string]any{}
		_ = json.Unmarshal([]byte(testIndividualDatastreamSnapshot), &data)
		reply = map[string]interface{}{"data": data}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/devices/%s/interfaces/%s/an/endpoint", testRealmName, testDeviceID, testServerOwnedInterfaceName):
		// receive data(stream)
		reply = map[string]interface{}{"data": ""}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/devices/%s/interfaces/%s/other/endpoint", testRealmName, testDeviceID, testServerOwnedInterfaceName):
		// receive data(stream)
		reply = map[string]interface{}{"data": ""}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/devices/%s/interfaces/%s/an/endpoint", testRealmName, testDeviceID, testServerOwnedPropertyInterfaceName):
		if req.Method == http.MethodPut {
			// set property
			reply = map[string]interface{}{"data": ""}
		} else if req.Method == http.MethodDelete {
			// unset property
			reply = map[string]interface{}{"data": ""}
			w.WriteHeader(http.StatusNoContent)
		}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/groups", testRealmName):
		// create group
		payload := DevicesAndGroup{Devices: testDeviceIDs, GroupName: testGroupName}
		reply = map[string]interface{}{"data": payload}
		w.WriteHeader(http.StatusCreated)
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/groups/%s/devices", testRealmName, url.PathEscape(testGroupName)):
		if req.Method == http.MethodGet {
			// list devices in a group
			reply = map[string]interface{}{"data": testDeviceIDs, "links": testGroupLinks}
		} else if req.Method == http.MethodPost {
			// add device to group
			reply = map[string]interface{}{"data": ""}
			w.WriteHeader(http.StatusCreated)
		}
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/groups/%s/devices/%s", testRealmName, url.PathEscape(testGroupName), testDeviceID):
		// remove device from group
		reply = map[string]interface{}{"data": ""}
		w.WriteHeader(http.StatusNoContent)
	}
	json.NewEncoder(w).Encode(reply)
}

func getTestContext(t *testing.T) (*Client, *httptest.Server) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(astarteAPIMock))

	// Use Client & URL from our local test server
	client, err := New(
		WithBaseURL(server.URL),
		WithJWT(testTokenValue),
		WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}

	return client, server
}
