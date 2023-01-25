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

func TestListDevices(t *testing.T) {
	c, _ := getTestContext(t)
	paginator, err := c.GetDeviceListPaginator(testRealmName, 10, DeviceIDFormat)
	if err != nil {
		t.Fatal(err)
	}
	if !paginator.HasNextPage() {
		t.Error("Paginator should have next page")
	}
	nextPageCall, err := paginator.GetNextPage()
	if err != nil {
		t.Fatal(err)
	}
	res, err := nextPageCall.Run(c)
	if err != nil {
		t.Fatal(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Fatal(err)
	}
	response, ok := data.([]string)
	if !ok {
		t.Fatal("Could not cast data correctly")
	}
	for i := 0; i < len(response); i++ {
		if response[i] != testDeviceIDs[i] {
			t.Fatalf("Different vaues when retrieving device IDs: %s vs %s", response[i], testDeviceIDs[i])
		}
	}
	if paginator.HasNextPage() {
		t.Error("Paginator should NOT have next page")
	}
	if _, err = paginator.GetNextPage(); err == nil {
		t.Error("Paginator should NOT have next page")
	}
}
