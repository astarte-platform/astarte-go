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

func TestListGroupDevices(t *testing.T) {
	c, _ := getTestContext(t)
	paginator, err := c.ListGroupDevices(testRealmName, testGroupName, 10, DeviceIDFormat)
	if err != nil {
		t.Error(err)
	}
	if !paginator.HasNextPage() {
		t.Error("Paginator should have next page")
	}
	nextPageCall, err := paginator.GetNextPage()
	if err != nil {
		t.Error(err)
	}
	res, err := nextPageCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	response, ok := data.([]string)
	if !ok {
		t.Error("Could not cast data correctly")
	}
	for i := 0; i < len(response); i++ {
		if response[i] != testDeviceIDs[i] {
			t.Errorf("Different vaues when retrieving device IDs: %s vs %s", response[i], testDeviceIDs[i])
		}
	}
	if paginator.HasNextPage() {
		t.Error("Paginator should NOT have next page")
	}
	if _, err = paginator.GetNextPage(); err == nil {
		t.Error("Paginator should NOT have next page")
	}
}

func TestCreateGroup(t *testing.T) {
	c, _ := getTestContext(t)
	createGroupRequest, err := c.CreateGroup(testRealmName, testGroupName, testDeviceIDs)
	if err != nil {
		t.Error(err)
	}
	res, err := createGroupRequest.Run(c)
	if err != nil {
		t.Error(err)
	}
	response, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	data, ok := response.(DevicesAndGroup)
	if !ok {
		t.Errorf("Expected a DevicesAndGroupPayload, found %v of type %T", response, response)
	}
	if data.GroupName != testGroupName {
		t.Errorf("Found unexpected group name: %s", data.GroupName)
	}
	for i := 0; i < len(testDeviceIDs); i++ {
		if data.Devices[i] != testDeviceIDs[i] {
			t.Errorf("Different vaues when retrieving device IDs: %s vs %s", data.Devices[i], testDeviceIDs[i])
		}
	}
}

func TestAddDeviceToGroup(t *testing.T) {
	c, _ := getTestContext(t)
	addDeviceToGroupCall, err := c.AddDeviceToGroup(testRealmName, testGroupName, testDeviceID)
	if err != nil {
		t.Error(err)
	}
	res, err := addDeviceToGroupCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveDeviceFromGroup(t *testing.T) {
	c, _ := getTestContext(t)
	removeDeviceFromGroupCall, err := c.RemoveDeviceFromGroup(testRealmName, testGroupName, testDeviceID)
	if err != nil {
		t.Error(err)
	}
	res, err := removeDeviceFromGroupCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}
