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
	"strings"
	"testing"

	"github.com/astarte-platform/astarte-go/interfaces"
)

func TestListInterfaces(t *testing.T) {
	c, _ := getTestContext(t)
	listInterfacesCall, err := c.ListInterfaces(testRealmName)
	if err != nil {
		t.Fatal(err)
	}
	res, err := listInterfacesCall.Run(c)
	if err != nil {
		t.Fatal(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Fatal(err)
	}
	interfaces, _ := data.([]string)
	for i := 0; i < len(testInterfacesList); i++ {
		if interfaces[i] != testInterfacesList[i] {
			t.Errorf("Listed interfaces not matching: %s vs %s", interfaces[i], testInterfacesList[i])
		}
	}
}

func TestListInterfaceMajorVersions(t *testing.T) {
	c, _ := getTestContext(t)
	listInterfaceMajorVersionsCall, err := c.ListInterfaceMajorVersions(testRealmName, testInterfaceName)
	if err != nil {
		t.Fatal(err)
	}
	res, err := listInterfaceMajorVersionsCall.Run(c)
	if err != nil {
		t.Fatal(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Fatal(err)
	}
	majors, _ := data.([]int)
	for i := 0; i < len(testInterfaceMajors); i++ {
		if majors[i] != testInterfaceMajors[i] {
			t.Errorf("Listed interface majors not matching: %d vs %d", majors[i], testInterfaceMajors[i])
		}
	}
}

func TestGetInterface(t *testing.T) {
	c, _ := getTestContext(t)
	createRealmCall, err := c.GetInterface(testRealmName, testInterfaceName, testInterfaceMajor)
	if err != nil {
		t.Error(err)
	}
	res, err := createRealmCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	testIface, _ := interfaces.ParseInterface([]byte(testInterface))
	iface, _ := data.(interfaces.AstarteInterface)
	//let's just assume it's enough
	if iface.Name != testIface.Name || iface.MajorVersion != testIface.MajorVersion || iface.MinorVersion != testIface.MinorVersion || iface.Type != testIface.Type {
		t.Error("Failed getting interface, different interface values")
	}
}

func TestInstallInterface(t *testing.T) {
	testIface, _ := interfaces.ParseInterface([]byte(testInterface))
	fmt.Println(testIface)

	c, _ := getTestContext(t)
	installInterfaceCall, err := c.InstallInterface(testRealmName, testIface, true)
	if err != nil {
		t.Error(err)
	}
	res, err := installInterfaceCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	iface, _ := data.(interfaces.AstarteInterface)
	fmt.Println(iface)
	//let's just assume it's enough
	if iface.Name != testIface.Name || iface.MajorVersion != testIface.MajorVersion || iface.MinorVersion != testIface.MinorVersion || iface.Type != testIface.Type {
		t.Error("Failed installing interface, different interface values")
	}
}

func TestInstallInterface_URL(t *testing.T) {
	testIface, _ := interfaces.ParseInterface([]byte(testInterface))
	c, _ := getTestContext(t)

	tests := []struct {
		name    string
		isAsync bool
		url     string
	}{
		{
			name:    "Async operation",
			isAsync: true,
			url:     "/realmmanagement/v1/test/interfaces'",
		},
		{
			name:    "Sync operation",
			isAsync: false,
			url:     "/realmmanagement/v1/test/interfaces?async_operation=false'",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			installInterfaceCall, err := c.InstallInterface(testRealmName, testIface, tc.isAsync)
			if err != nil {
				t.Error(err)
			}

			if !strings.Contains(installInterfaceCall.ToCurl(c), tc.url) {
				t.Errorf("%v does not contain %s", installInterfaceCall.ToCurl(c), tc.url)
			}
		})
	}
}

func TestDeleteInterface(t *testing.T) {
	c, _ := getTestContext(t)
	deleteInterfaceCall, err := c.DeleteInterface(testRealmName, testInterfaceName, testInterfaceMajor)
	if err != nil {
		t.Error(err)
	}
	res, err := deleteInterfaceCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateInterface(t *testing.T) {
	testIface, _ := interfaces.ParseInterface([]byte(testInterface))
	c, _ := getTestContext(t)
	updateInterfaceCall, err := c.UpdateInterface(testRealmName, testInterfaceName, testInterfaceMajor, testIface, true)
	if err != nil {
		t.Error(err)
	}
	res, err := updateInterfaceCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateInterface_URL(t *testing.T) {
	testIface, _ := interfaces.ParseInterface([]byte(testInterface))
	c, _ := getTestContext(t)

	tests := []struct {
		name    string
		isAsync bool
		url     string
	}{
		{
			name:    "Async operation",
			isAsync: true,
			url:     "/realmmanagement/v1/test/interfaces/ah.yes.an.Interface/1'",
		},
		{
			name:    "Sync operation",
			isAsync: false,
			url:     "/realmmanagement/v1/test/interfaces/ah.yes.an.Interface/1?async_operation=false'",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			updateInterfaceCall, err := c.UpdateInterface(
				testRealmName,
				testInterfaceName,
				testInterfaceMajor,
				testIface,
				tc.isAsync,
			)
			if err != nil {
				t.Error(err)
			}

			if !strings.Contains(updateInterfaceCall.ToCurl(c), tc.url) {
				t.Errorf("%v does not contain %s", updateInterfaceCall.ToCurl(c), tc.url)
			}
		})
	}
}

func TestListTriggers(t *testing.T) {
	c, _ := getTestContext(t)
	listTriggersCall, err := c.ListTriggers(testRealmName)
	if err != nil {
		t.Fatal(err)
	}
	res, err := listTriggersCall.Run(c)
	if err != nil {
		t.Fatal(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Fatal(err)
	}
	triggers, _ := data.([]string)
	for i := 0; i < len(testTriggersList); i++ {
		if triggers[i] != testTriggersList[i] {
			t.Errorf("Listed interfaces not matching: %s vs %s", triggers[i], testTriggersList[i])
		}
	}
}

func TestGetTrigger(t *testing.T) {
	c, _ := getTestContext(t)
	getTriggerCall, err := c.GetTrigger(testRealmName, testTriggerName)
	if err != nil {
		t.Error(err)
	}
	res, err := getTriggerCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	triggerMap, _ := data.(map[string]interface{})
	triggerName, _ := triggerMap["name"].(string)

	//let's just assume it's enough
	if triggerName != testTriggerName {
		t.Error("Failed getting trigger, different trigger values")
	}
}

func TestInstallTrigger(t *testing.T) {
	c, _ := getTestContext(t)
	trigger := map[string]any{}
	_ = json.Unmarshal([]byte(testTrigger), &trigger)
	installTriggerCall, err := c.InstallTrigger(testRealmName, trigger)
	if err != nil {
		t.Error(err)
	}
	res, err := installTriggerCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	triggerMap, _ := data.(map[string]interface{})
	triggerName, _ := triggerMap["name"].(string)

	//let's just assume it's enough
	if triggerName != testTriggerName {
		t.Error("Failed getting trigger, different trigger values")
	}
}

func TestDeleteTrigger(t *testing.T) {
	c, _ := getTestContext(t)
	deleteTriggerCall, err := c.DeleteTrigger(testRealmName, testTriggerName)
	if err != nil {
		t.Error(err)
	}
	res, err := deleteTriggerCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}

func TestListTriggerDeliveryPolicies(t *testing.T) {
	c, _ := getTestContext(t)
	listPoliciesCall, err := c.ListTriggerDeliveryPolicies(testRealmName)
	if err != nil {
		t.Fatal(err)
	}
	res, err := listPoliciesCall.Run(c)
	if err != nil {
		t.Fatal(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Fatal(err)
	}
	policies, _ := data.([]string)
	for i := 0; i < len(testPoliciesList); i++ {
		if policies[i] != testPoliciesList[i] {
			t.Errorf("Listed policies not matching: %s vs %s", policies[i], testPoliciesList[i])
		}
	}
}

func TestGetTriggerDeliveryPolicy(t *testing.T) {
	c, _ := getTestContext(t)
	getPolicyCall, err := c.GetTriggerDeliveryPolicy(testRealmName, testPolicyName)
	if err != nil {
		t.Error(err)
	}
	res, err := getPolicyCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	policyMap, _ := data.(map[string]interface{})
	policyName, _ := policyMap["name"].(string)

	//let's just assume it's enough
	if policyName != testPolicyName {
		t.Error("Failed getting policy, different policy values")
	}
}

func TestInstallTriggerDeliveryPolicy(t *testing.T) {
	c, _ := getTestContext(t)
	policy := map[string]any{}
	_ = json.Unmarshal([]byte(testPolicy), &policy)
	installPolicyCall, err := c.InstallTriggerDeliveryPolicy(testRealmName, policy)
	if err != nil {
		t.Error(err)
	}
	res, err := installPolicyCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	policyMap, _ := data.(map[string]interface{})
	policyName, _ := policyMap["name"].(string)

	//let's just assume it's enough
	if policyName != testPolicyName {
		t.Error("Failed getting policy, different policy values")
	}
}

func TestDeleteTriggerDeliveryPolicy(t *testing.T) {
	c, _ := getTestContext(t)
	deletePolicyCall, err := c.DeleteTriggerDeliveryPolicy(testRealmName, testPolicyName)
	if err != nil {
		t.Error(err)
	}
	res, err := deletePolicyCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}
