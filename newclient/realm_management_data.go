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
	"encoding/json"
	"io"

	"github.com/astarte-platform/astarte-go/interfaces"
	"github.com/tidwall/gjson"
)

// Parses data obtained by performing a request to list interfaces in a realm.
// Returns the list of interface names as an array of strings.
func (r listInterfacesResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	ret := []string{}
	for _, v := range gjson.GetBytes(b, "data").Array() {
		ret = append(ret, v.Str)
	}
	return ret, nil
}
func (r listInterfacesResponse) Raw() {}

// Parses data obtained by performing a request to list an interface's major versions.
// Returns the list of versions as an array of ints.
func (r listInterfaceMajorVersionsResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	ret := []int{}
	for _, v := range gjson.GetBytes(b, "data").Array() {
		ret = append(ret, int(v.Num))
	}
	return ret, nil
}
func (r listInterfaceMajorVersionsResponse) Raw() {}

// Parses data obtained by performing a request to retrieve an interface.
// Returns the interface as an interfaces.AstarteInterface.
func (r getInterfaceResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := interfaces.AstarteInterface{}
	// TODO check err
	_ = json.Unmarshal(v, &ret)
	return interfaces.EnsureInterfaceDefaults(ret), nil

}
func (r getInterfaceResponse) Raw() {}

// Parses data obtained by performing a request to install an interface.
// Returns the interface as an interfaces.AstarteInterface.
func (r installInterfaceResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := interfaces.AstarteInterface{}
	// TODO check err
	_ = json.Unmarshal(v, &ret)
	return interfaces.EnsureInterfaceDefaults(ret), nil
}

func (r installInterfaceResponse) Raw() {}

// Parses data obtained by performing a request to delete an interface.
// The returned values do not matter.
func (r deleteInterfaceResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	return "", nil
}

func (r deleteInterfaceResponse) Raw() {}

// Parses data obtained by performing a request to update an interface.
// The returned values do not matter.
func (r updateInterfaceResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	return "", nil
}

func (r updateInterfaceResponse) Raw() {}

// Parses data obtained by performing a request to list triggers in a realm.
// Returns the list of triggers names as an array of strings.
func (r listTriggersResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	ret := []string{}
	for _, v := range gjson.GetBytes(b, "data").Array() {
		ret = append(ret, v.Str)
	}
	return ret, nil
}
func (r listTriggersResponse) Raw() {}

// Parses data obtained by performing a request to retrieve a trigger.
// Returns the trigger payload as a map[string]any.
func (r getTriggerResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := map[string]any{}
	err := json.Unmarshal(v, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r getTriggerResponse) Raw() {}

// Parses data obtained by performing a request to install a trigger.
// Returns the trigger payload as a map[string]any.
func (r installTriggerResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := map[string]any{}
	err := json.Unmarshal(v, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r installTriggerResponse) Raw() {}

// Parses data obtained by performing a request to delete a trigger.
// The returned values do not matter.
func (r deleteTriggerResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	return "", nil
}

func (r deleteTriggerResponse) Raw() {}
