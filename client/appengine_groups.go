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
	"fmt"
	"net/url"
	"path"
)

// This file contains all API Calls related to device group management

// ListGroups lists the groups in a Realm
func (s *AppEngineService) ListGroups(realm string, token string) ([]string, error) {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups", realm))
	decoder, err := s.client.genericJSONDataAPIGET(callURL.String(), token, 200)
	if err != nil {
		return nil, err
	}
	var responseBody struct {
		Data []string `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Data, nil
}

// CreateGroup creates a group with the given deviceIdentifierList in the Realm
func (s *AppEngineService) CreateGroup(realm string, groupName string, deviceIdentifierList []string,
	deviceIdentifiersType DeviceIdentifierType, token string) error {

	deviceIDList := make([]string, len(deviceIdentifierList))
	for i, deviceIdentifier := range deviceIdentifierList {
		deviceID, err := s.GetDeviceIDFromDeviceIdentifier(realm, deviceIdentifier, deviceIdentifiersType, token)
		if err != nil {
			return err
		}
		deviceIDList[i] = deviceID
	}
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups", realm))
	payload := map[string]interface{}{"group_name": groupName, "devices": deviceIDList}
	err := s.client.genericJSONDataAPIPost(callURL.String(), payload, token, 201)
	if err != nil {
		return err
	}

	return nil
}

// ListGroupDevices lists the devices that belong to a group
func (s *AppEngineService) ListGroupDevices(realm string, groupName string, token string) ([]string, error) {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups/%s/devices", realm, url.PathEscape(groupName)))
	decoder, err := s.client.genericJSONDataAPIGET(callURL.String(), token, 200)
	if err != nil {
		return nil, err
	}
	var responseBody struct {
		Data []string `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Data, nil
}

// AddDeviceToGroup adds a device to the group
func (s *AppEngineService) AddDeviceToGroup(realm string, groupName string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType, token string) error {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups/%s/devices", realm, url.PathEscape(groupName)))
	deviceID, err := s.GetDeviceIDFromDeviceIdentifier(realm, deviceIdentifier, deviceIdentifierType, token)
	if err != nil {
		return err
	}
	payload := map[string]string{"device_id": deviceID}
	err = s.client.genericJSONDataAPIPost(callURL.String(), payload, token, 201)
	if err != nil {
		return err
	}

	return nil
}

// RemoveDeviceFromGroup removes a device from the group
func (s *AppEngineService) RemoveDeviceFromGroup(realm string, groupName string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType, token string) error {
	deviceID, err := s.GetDeviceIDFromDeviceIdentifier(realm, deviceIdentifier, deviceIdentifierType, token)
	if err != nil {
		return err
	}
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups/%s/devices/%s", realm, url.PathEscape(groupName), deviceID))
	err = s.client.genericJSONDataAPIDelete(callURL.String(), token, 204)
	if err != nil {
		return err
	}

	return nil
}
