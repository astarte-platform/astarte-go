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
	"net/http"

	"github.com/astarte-platform/astarte-go/interfaces"
	"github.com/tidwall/gjson"
)

// Parses data obtained by performing a request to list interfaces in a realm.
// Returns the list of interface names as an array of strings.
func (r ListInterfacesResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	ret := []string{}
	for _, v := range gjson.GetBytes(b, "data").Array() {
		ret = append(ret, v.Str)
	}
	return ret, nil
}
func (r ListInterfacesResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to list an interface's major versions.
// Returns the list of versions as an array of ints.
func (r ListInterfaceMajorVersionsResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	ret := []int{}
	for _, v := range gjson.GetBytes(b, "data").Array() {
		ret = append(ret, int(v.Num))
	}
	return ret, nil
}
func (r ListInterfaceMajorVersionsResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to retrieve an interface.
// Returns the interface as an interfaces.AstarteInterface.
func (r GetInterfaceResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := interfaces.AstarteInterface{}
	// TODO check err
	_ = json.Unmarshal(v, &ret)
	return interfaces.EnsureInterfaceDefaults(ret), nil

}
func (r GetInterfaceResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to install an interface.
// Returns the interface as an interfaces.AstarteInterface.
func (r InstallInterfaceResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := interfaces.AstarteInterface{}
	// TODO check err
	_ = json.Unmarshal(v, &ret)
	return interfaces.EnsureInterfaceDefaults(ret), nil
}

func (r InstallInterfaceResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to list triggers in a realm.
// Returns the list of triggers names as an array of strings.
func (r ListTriggersResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	ret := []string{}
	for _, v := range gjson.GetBytes(b, "data").Array() {
		ret = append(ret, v.Str)
	}
	return ret, nil
}
func (r ListTriggersResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to retrieve a trigger.
// Returns the trigger payload as a map[string]any.
func (r GetTriggerResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := map[string]any{}
	err := json.Unmarshal(v, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r GetTriggerResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to install a trigger.
// Returns the trigger payload as a map[string]any.
func (r InstallTriggerResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	v := []byte(gjson.GetBytes(b, "data").Raw)
	ret := map[string]any{}
	err := json.Unmarshal(v, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r InstallTriggerResponse) Raw() *http.Response {
	return r.res
}
