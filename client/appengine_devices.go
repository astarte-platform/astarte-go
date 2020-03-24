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

// This file contains all API Calls related to device management and information such as aliases, stats...

// DeviceIdentifierType represents what kind of identifier is used for identifying a Device.
type DeviceIdentifierType int

const (
	// AutodiscoverDeviceIdentifier is the default, and uses heuristics to autodetermine which kind of
	// identifier is being used.
	AutodiscoverDeviceIdentifier DeviceIdentifierType = iota
	// AstarteDeviceID is the Device's ID in its standard format.
	AstarteDeviceID
	// AstarteDeviceAlias is one of the Device's Aliases.
	AstarteDeviceAlias
)

// ListDevices returns a list of Devices in the Realm
func (s *AppEngineService) ListDevices(realm string, token string) ([]string, error) {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices", realm))
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

// GetDevice returns the DeviceDetails of a single Device in the Realm
func (s *AppEngineService) GetDevice(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, token string) (DeviceDetails, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	decoder, err := s.client.genericJSONDataAPIGET(callURL.String(), token, 200)
	if err != nil {
		return DeviceDetails{}, err
	}
	var responseBody struct {
		Data DeviceDetails `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return DeviceDetails{}, err
	}

	return responseBody.Data, nil
}

// GetDeviceIDFromDeviceIdentifier returns the DeviceID of a Device identified with a deviceIdentifier
// of type deviceIdentifierType.
func (s *AppEngineService) GetDeviceIDFromDeviceIdentifier(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType, token string) (string, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	switch resolvedDeviceIdentifierType {
	case AstarteDeviceAlias:
		return s.GetDeviceIDFromAlias(realm, deviceIdentifier, token)
	default:
		return deviceIdentifier, nil
	}
}

// GetDeviceIDFromAlias returns the Device ID of a device given one of its aliases
func (s *AppEngineService) GetDeviceIDFromAlias(realm string, deviceAlias string, token string) (string, error) {
	deviceDetails, err := s.GetDevice(realm, deviceAlias, AstarteDeviceAlias, token)
	if err != nil {
		return "", err
	}

	return deviceDetails.DeviceID, nil
}

// ListDeviceInterfaces returns the list of Interfaces exposed by the Device's introspection
func (s *AppEngineService) ListDeviceInterfaces(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType, token string) ([]string, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/interfaces", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
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

// ListDeviceAliases is an helper to list all aliases of a Device
func (s *AppEngineService) ListDeviceAliases(realm string, deviceID string, token string) (map[string]string, error) {
	deviceDetails, err := s.GetDevice(realm, deviceID, AstarteDeviceID, token)
	if err != nil {
		return nil, err
	}
	return deviceDetails.Aliases, nil
}

// AddDeviceAlias adds an Alias to a Device
func (s *AppEngineService) AddDeviceAlias(realm string, deviceID string, aliasTag string, deviceAlias string, token string) error {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s", realm, deviceID))
	payload := map[string]map[string]string{"aliases": {aliasTag: deviceAlias}}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, token, 200)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDeviceAlias deletes an Alias from a Device based on the Alias' tag
func (s *AppEngineService) DeleteDeviceAlias(realm string, deviceID string, aliasTag string, token string) error {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s", realm, deviceID))
	// We're using map[string]interface{} rather than map[string]string since we want to have null
	// rather than an empty string in the JSON payload, and this is the only way.
	payload := map[string]map[string]interface{}{"aliases": {aliasTag: nil}}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, token, 200)
	if err != nil {
		return err
	}

	return nil
}

// InhibitDevice sets the Credentials Inhibition state of a Device
func (s *AppEngineService) InhibitDevice(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType, token string, inhibit bool) error {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	payload := map[string]bool{"credentials_inhibited": inhibit}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, token, 200)
	if err != nil {
		return err
	}

	return nil
}

// GetDevicesStats returns the DevicesStats of a Realm
func (s *AppEngineService) GetDevicesStats(realm string, token string) (DevicesStats, error) {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/stats/devices", realm))
	decoder, err := s.client.genericJSONDataAPIGET(callURL.String(), token, 200)
	if err != nil {
		return DevicesStats{}, err
	}
	var responseBody struct {
		Data DevicesStats `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return DevicesStats{}, err
	}

	return responseBody.Data, nil
}
