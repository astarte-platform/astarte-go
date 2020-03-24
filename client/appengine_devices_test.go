// Copyright © 2019 Ispirata Srl
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
	"reflect"
	"testing"
)

func TestListDevices(t *testing.T) {
	// Start a local HTTP server
	client, server := getTestContext(t)
	// Close the server when test finishes
	defer server.Close()

	devices, err := client.AppEngine.ListDevices(testRealmName)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(devices, testDevices) {
		t.Log(devices)
		t.Log(testDevices)
		t.Fail()
	}
}
