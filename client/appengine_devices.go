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

// ListDevices returns the list of Device IDs for all Devices in the Realm. The
// returned result can be large, GetDeviceListPaginator can be used instead to
// retrieve the device list incrementally.
func (s *AppEngineService) ListDevices(realm string) ([]string, error) {
	result := []string{}

	paginator, err := s.GetDeviceListPaginator(realm, defaultPageSize, DeviceIDFormat)
	if err != nil {
		return result, err
	}

	for hasNext := paginator.HasNextPage(); hasNext; hasNext = paginator.HasNextPage() {
		page := []string{}
		err := paginator.GetNextPage(&page)
		if err != nil {
			return []string{}, err
		}
		result = append(result, page...)
	}

	return result, nil
}

// ListDevicesWithDetails returns a list of all Devices in the Realm, each
// represented by a DeviceDetails struct. The returned result can be large,
// GetDeviceListPaginator can be used instead to retrieve the device list
// incrementally.
func (s *AppEngineService) ListDevicesWithDetails(realm string) ([]DeviceDetails, error) {
	result := []DeviceDetails{}

	paginator, err := s.GetDeviceListPaginator(realm, defaultPageSize, DeviceDetailsFormat)
	if err != nil {
		return result, err
	}

	for hasNext := paginator.HasNextPage(); hasNext; hasNext = paginator.HasNextPage() {
		page := []DeviceDetails{}
		err := paginator.GetNextPage(&page)
		if err != nil {
			return []DeviceDetails{}, err
		}
		result = append(result, page...)
	}

	return result, nil
}

// GetDeviceListPaginator returns a Paginator for all the Devices in the realm.
// The paginator can return different result formats depending on the format
// parameter.
func (s *AppEngineService) GetDeviceListPaginator(realm string, pageSize int, format DeviceResultFormat) (DeviceListPaginator, error) {
	callURL, err := url.Parse(s.appEngineURL.String())
	if err != nil {
		return DeviceListPaginator{}, err
	}
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices", realm))
	query := url.Values{}

	deviceListPaginator := DeviceListPaginator{
		baseURL:     callURL,
		nextQuery:   query,
		format:      format,
		pageSize:    pageSize,
		client:      s.client,
		hasNextPage: true,
	}
	return deviceListPaginator, nil
}

// GetDevice returns the DeviceDetails of a single Device in the Realm
func (s *AppEngineService) GetDevice(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) (DeviceDetails, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	deviceDetails := DeviceDetails{}
	err := s.client.genericJSONDataAPIGET(&deviceDetails, callURL.String(), 200)

	return deviceDetails, err
}

// GetDeviceIDFromDeviceIdentifier returns the DeviceID of a Device identified with a deviceIdentifier
// of type deviceIdentifierType.
func (s *AppEngineService) GetDeviceIDFromDeviceIdentifier(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType) (string, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	switch resolvedDeviceIdentifierType {
	case AstarteDeviceAlias:
		return s.GetDeviceIDFromAlias(realm, deviceIdentifier)
	default:
		return deviceIdentifier, nil
	}
}

// GetDeviceIDFromAlias returns the Device ID of a device given one of its aliases
func (s *AppEngineService) GetDeviceIDFromAlias(realm string, deviceAlias string) (string, error) {
	deviceDetails, err := s.GetDevice(realm, deviceAlias, AstarteDeviceAlias)
	if err != nil {
		return "", err
	}

	return deviceDetails.DeviceID, nil
}

// ListDeviceInterfaces returns the list of Interfaces exposed by the Device's introspection
func (s *AppEngineService) ListDeviceInterfaces(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType) ([]string, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/interfaces", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	deviceInterfacesList := []string{}
	err := s.client.genericJSONDataAPIGET(&deviceInterfacesList, callURL.String(), 200)

	return deviceInterfacesList, err
}

// ListDeviceAliases is an helper to list all aliases of a Device
func (s *AppEngineService) ListDeviceAliases(realm string, deviceID string) (map[string]string, error) {
	deviceDetails, err := s.GetDevice(realm, deviceID, AstarteDeviceID)
	if err != nil {
		return nil, err
	}
	return deviceDetails.Aliases, nil
}

// AddDeviceAlias adds an Alias to a Device
func (s *AppEngineService) AddDeviceAlias(realm string, deviceID string, aliasTag string, deviceAlias string) error {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s", realm, deviceID))
	payload := map[string]map[string]string{"aliases": {aliasTag: deviceAlias}}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, 200)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDeviceAlias deletes an Alias from a Device based on the Alias' tag
func (s *AppEngineService) DeleteDeviceAlias(realm string, deviceID string, aliasTag string) error {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s", realm, deviceID))
	// We're using map[string]interface{} rather than map[string]string since we want to have null
	// rather than an empty string in the JSON payload, and this is the only way.
	payload := map[string]map[string]interface{}{"aliases": {aliasTag: nil}}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, 200)
	if err != nil {
		return err
	}

	return nil
}

// InhibitDevice sets the Credentials Inhibition state of a Device
func (s *AppEngineService) InhibitDevice(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType, inhibit bool) error {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	payload := map[string]bool{"credentials_inhibited": inhibit}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, 200)
	if err != nil {
		return err
	}

	return nil
}

// GetDevicesStats returns the DevicesStats of a Realm
func (s *AppEngineService) GetDevicesStats(realm string) (DevicesStats, error) {
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/stats/devices", realm))
	deviceStats := DevicesStats{}
	err := s.client.genericJSONDataAPIGET(&deviceStats, callURL.String(), 200)

	return deviceStats, err
}

// ListDeviceAttributes is an helper to list all Attributes of a Device
func (s *AppEngineService) ListDeviceAttributes(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) (map[string]string, error) {
	deviceDetails, err := s.GetDevice(realm, deviceIdentifier, deviceIdentifierType)
	if err != nil {
		return nil, err
	}
	return deviceDetails.Attributes, nil
}

// SetDeviceAttribute sets an Attribute key to a certain value for a Device
func (s *AppEngineService) SetDeviceAttribute(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, attributeKey, attributeValue string) error {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	payload := map[string]map[string]string{"attributes": {attributeKey: attributeValue}}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, 200)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDeviceAttribute deletes an Attribute key and its value from a Device
func (s *AppEngineService) DeleteDeviceAttribute(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, attributeKey string) error {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(s.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType)))
	// We're using map[string]interface{} rather than map[string]string since we want to have null
	// rather than an empty string in the JSON payload, and this is the only way.
	payload := map[string]map[string]interface{}{"attributes": {attributeKey: nil}}
	err := s.client.genericJSONDataAPIPatch(callURL.String(), payload, 200)
	if err != nil {
		return err
	}

	return nil
}
