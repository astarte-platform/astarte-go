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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testRealmName  = "test"
	testTokenValue = "bogus"
)

var testInterfaces map[string]string = map[string]string{
	"org.astarte-platform.genericsensors.Values": `{
		"interface_name": "org.astarte-platform.genericsensors.Values",
		"version_major": 0,
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
	}`,
	"org.astarte-platform.genericsensors.SamplingRate": `{
		"interface_name": "org.astarte-platform.genericsensors.SamplingRate",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "server",
		"description": "Configure sensors sampling rate and enable/disable.",
		"doc": "This interface allows to set generic sensors sampling rate and enable/disable policies for each sensor. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/enable",
				"type": "boolean",
				"allow_unset": true,
				"description": "Enable/disable sensor data transmission.",
				"doc": "When true sampled data transmission is always on, otherwise when false is off. When unset data transmission policy is up to the sensor."
			},
			{
				"endpoint": "/%{sensor_id}/samplingPeriod",
				"type": "integer",
				"allow_unset": true,
				"description": "Sensor sample transmission period.",
				"doc": "Send a sampled value every samplingPeriod seconds. When unset sampling period is up to the sensor."
			}
		]
	}`,
	"org.astarte-platform.genericsensors.AvailableSensors": `{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/name",
				"type": "string",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name."
			},
			{
				"endpoint": "/%{sensor_id}/unit",
				"type": "string",
				"description": "Sample data measurement unit.",
				"doc": "SI unit such as m, kg, K, etc..."
			}
		]
	}`,
}

var testDevices []string = []string{"1vMeFtaJQF259nMsnis3sw", "t1J1uQSBQRi_1F3zIrjyYw", "V_pY-ZrLQzWz4iGjGu-NuQ"}

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
	switch {
	case req.URL.Path == fmt.Sprintf("/appengine/v1/%s/devices", testRealmName):
		links := map[string]string{"self": fmt.Sprintf("/v1/%s/devices", testRealmName)}
		reply := map[string]interface{}{"data": testDevices, "links": links}
		json.NewEncoder(w).Encode(reply)
	}
}

func getTestContext(t *testing.T) (*Client, *httptest.Server) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(astarteAPIMock))

	// Use Client & URL from our local test server
	client, err := NewClient(server.URL, server.Client())
	if err != nil {
		t.Error(err)
	}
	client.SetToken(testTokenValue)

	return client, server
}
