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
	"fmt"
	"net/url"
	"path"

	"github.com/astarte-platform/astarte-go/interfaces"
)

// RealmManagementService is the API Client for RealmManagement API
type RealmManagementService struct {
	client             *Client
	realmManagementURL *url.URL
}

// ListInterfaces returns all interfaces in a Realm.
func (s *RealmManagementService) ListInterfaces(realm string) ([]string, error) {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces", realm))

	interfacesList := []string{}
	err := s.client.genericJSONDataAPIGET(&interfacesList, callURL.String(), 200)

	return interfacesList, err
}

// ListInterfaceMajorVersions returns all available major versions for a given Interface in a Realm.
func (s *RealmManagementService) ListInterfaceMajorVersions(realm string, interfaceName string) ([]int, error) {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s", realm, interfaceName))

	interfaceMajorVersions := []int{}
	err := s.client.genericJSONDataAPIGET(&interfaceMajorVersions, callURL.String(), 200)

	return interfaceMajorVersions, err
}

// GetInterface returns an interface, identified by a Major version, in a Realm
func (s *RealmManagementService) GetInterface(realm string, interfaceName string, interfaceMajor int) (interfaces.AstarteInterface, error) {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s/%v", realm, interfaceName, interfaceMajor))

	iface := interfaces.AstarteInterface{}
	err := s.client.genericJSONDataAPIGET(&iface, callURL.String(), 200)

	return interfaces.EnsureInterfaceDefaults(iface), err
}

// InstallInterface installs a new major version of an Interface into the Realm
func (s *RealmManagementService) InstallInterface(realm string, interfacePayload interfaces.AstarteInterface) error {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces", realm))
	return s.client.genericJSONDataAPIPost(callURL.String(), interfacePayload, 201)
}

// DeleteInterface deletes a draft Interface from the Realm
func (s *RealmManagementService) DeleteInterface(realm string, interfaceName string, interfaceMajor int) error {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s/%v", realm, interfaceName, interfaceMajor))
	return s.client.genericJSONDataAPIDelete(callURL.String(), 204)
}

// UpdateInterface updates an existing major version of an Interface to a new minor.
func (s *RealmManagementService) UpdateInterface(realm string, interfaceName string, interfaceMajor int, interfacePayload interfaces.AstarteInterface) error {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s/%v", realm, interfaceName, interfaceMajor))
	return s.client.genericJSONDataAPIPut(callURL.String(), interfacePayload, 204)
}

// ListTriggers returns all triggers in a Realm.
func (s *RealmManagementService) ListTriggers(realm string) ([]string, error) {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers", realm))

	triggers := []string{}
	err := s.client.genericJSONDataAPIGET(&triggers, callURL.String(), 200)

	return triggers, err
}

// GetTrigger returns a trigger installed in a Realm
func (s *RealmManagementService) GetTrigger(realm string, triggerName string) (map[string]interface{}, error) {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers/%s", realm, triggerName))

	trigger := map[string]interface{}{}
	err := s.client.genericJSONDataAPIGET(&trigger, callURL.String(), 200)

	return trigger, err
}

// InstallTrigger installs a Trigger into the Realm
func (s *RealmManagementService) InstallTrigger(realm string, triggerPayload interface{}) error {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers", realm))
	return s.client.genericJSONDataAPIPost(callURL.String(), triggerPayload, 201)
}

// DeleteTrigger deletes a Trigger from the Realm
func (s *RealmManagementService) DeleteTrigger(realm string, triggerName string) error {
	callURL, _ := url.Parse(s.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers/%s", realm, triggerName))
	return s.client.genericJSONDataAPIDelete(callURL.String(), 204)
}
