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
	"fmt"
	"net/http"
	"net/url"

	"moul.io/http2curl"
)

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

// GetDeviceListPaginator returns a Paginator for all the Devices in the realm.
// The paginator can return different result formats depending on the format
// parameter.
func (c *Client) GetDeviceListPaginator(realm string, pageSize int, format DeviceResultFormat) (Paginator, error) {
	callURL := makeURL(c.appEngineURL, "/v1/%s/devices", realm)
	query := url.Values{}
	deviceListPaginator := DeviceListPaginator{
		baseURL:     callURL,
		nextQuery:   query,
		format:      format,
		pageSize:    pageSize,
		client:      c,
		hasNextPage: true,
	}

	return &deviceListPaginator, nil
}

type GetDeviceDetailsRequest struct {
	req     *http.Request
	expects int
}

// GetDevice builds a request to return the DeviceDetails of a single Device in the Realm.
func (c *Client) GetDeviceDetails(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL := makeURL(c.appEngineURL, "/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType))
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return GetDeviceDetailsRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r GetDeviceDetailsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return GetDeviceDetailsResponse{res: res}, nil
}

func (r GetDeviceDetailsRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type GetDeviceIDFromAliasRequest struct {
	req     *http.Request
	expects int
}

// GetDeviceIDFromAlias builds a request to return the Device ID of a device given one of its aliases.
func (c *Client) GetDeviceIDFromAlias(realm string, deviceAlias string) (AstarteRequest, error) {
	getDeviceRequest, err := c.GetDeviceDetails(realm, deviceAlias, AstarteDeviceAlias)
	if err != nil {
		return Empty{}, nil
	}
	getDeviceDetailsRequest, _ := getDeviceRequest.(GetDeviceDetailsRequest)
	return GetDeviceIDFromAliasRequest(getDeviceDetailsRequest), nil
}

// nolint:bodyclose
func (r GetDeviceIDFromAliasRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return GetDeviceIDFromAliasResponse{res: res}, nil
}

func (r GetDeviceIDFromAliasRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	// TODO check
	return fmt.Sprintf("%s | grep 'DeviceID'\n", command)
}

type ListDeviceInterfacesRequest struct {
	req     *http.Request
	expects int
}

// ListDeviceInterfaces builds a request to retrieve the list of interfaces exposed by the Device's introspection.
func (c *Client) ListDeviceInterfaces(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL := makeURL(c.appEngineURL, "/v1/%s/%s/interfaces", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType))
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return ListDeviceInterfacesRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r ListDeviceInterfacesRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListDeviceInterfacesResponse{res: res}, nil
}

func (r ListDeviceInterfacesRequest) ToCurl(_ *Client) string {
	return ""
}

type GetDevicesStatsRequest struct {
	req     *http.Request
	expects int
}

// GetDevicesStats builds a request to return the DevicesStats of a Realm.
func (c *Client) GetDevicesStats(realm string) (AstarteRequest, error) {
	callURL := makeURL(c.appEngineURL, "/v1/%s/stats/devices", realm)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return GetDevicesStatsRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r GetDevicesStatsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return GetDeviceStatsResponse{res: res}, nil
}

func (r GetDevicesStatsRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type ListDeviceAliasesRequest struct {
	req     *http.Request
	expects int
}

// ListDeviceAliases builds a request to list all aliases of a Device.
func (c *Client) ListDeviceAliases(realm string, deviceIdentifier string,
	deviceIdentifierType DeviceIdentifierType) (AstarteRequest, error) {
	getDeviceRequest, err := c.GetDeviceDetails(realm, deviceIdentifier, deviceIdentifierType)
	if err != nil {
		return Empty{}, nil
	}
	getDeviceDetailsRequest, _ := getDeviceRequest.(GetDeviceDetailsRequest)
	return ListDeviceAliasesRequest(getDeviceDetailsRequest), nil
}

// nolint:bodyclose
func (r ListDeviceAliasesRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListDeviceAliasesResponse{res: res}, nil
}

func (r ListDeviceAliasesRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	// TODO check
	return fmt.Sprintf("%s | grep 'Aliases'\n", command)
}

type AddDeviceAliasRequest struct {
	req     *http.Request
	expects int
}

// AddDeviceAlias builds a request to add an Alias to a Device
func (c *Client) AddDeviceAlias(realm string, deviceID string, aliasTag string, deviceAlias string) (AstarteRequest, error) {
	callURL := makeURL(c.appEngineURL, "/v1/%s/devices/%s", realm, deviceID)
	aliasMap := map[string]map[string]string{"aliases": {aliasTag: deviceAlias}}
	payload, _ := makeBody(aliasMap)
	req := c.makeHTTPrequestWithContentType(http.MethodPatch, callURL, payload, "application/merge-patch+json")

	return AddDeviceAliasRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r AddDeviceAliasRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}

	return NoDataResponse{res: res}, nil
}

func (r AddDeviceAliasRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type DeleteDeviceAliasRequest struct {
	req     *http.Request
	expects int
}

// DeleteDeviceAlias builds a request to delete an Alias from a Device based on the Alias' tag.
func (c *Client) DeleteDeviceAlias(realm string, deviceID string, aliasTag string) (AstarteRequest, error) {
	callURL := makeURL(c.appEngineURL, "/v1/%s/devices/%s", realm, deviceID)
	// We're using map[string]interface{} rather than map[string]string since we want to have null
	// rather than an empty string in the JSON payload, and this is the only way.
	aliasMap := map[string]map[string]interface{}{"aliases": {aliasTag: nil}}
	payload, _ := makeBody(aliasMap)
	req := c.makeHTTPrequestWithContentType(http.MethodPatch, callURL, payload, "application/merge-patch+json")

	return DeleteDeviceAliasRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r DeleteDeviceAliasRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r DeleteDeviceAliasRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	// TODO check
	return fmt.Sprint(command)
}

type InhibitDeviceRequest struct {
	req     *http.Request
	expects int
}

// SetDeviceInhibited builds a request to set the Credentials Inhibition state of a Device.
func (c *Client) SetDeviceInhibited(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, inhibit bool) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL := makeURL(c.appEngineURL, "/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType))
	credentialsMap := map[string]bool{"credentials_inhibited": inhibit}
	payload, _ := makeBody(credentialsMap)
	req := c.makeHTTPrequestWithContentType(http.MethodPatch, callURL, payload, "application/merge-patch+json")

	return InhibitDeviceRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r InhibitDeviceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	// no response expected
	return NoDataResponse{res: res}, nil
}

func (r InhibitDeviceRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	// TODO check
	return fmt.Sprint(command)
}

type ListDeviceAttributesRequest struct {
	req     *http.Request
	expects int
}

// ListDeviceAttributes builds a request to list all Attributes of a Device.
func (c *Client) ListDeviceAttributes(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) (AstarteRequest, error) {
	getDeviceRequest, err := c.GetDeviceDetails(realm, deviceIdentifier, deviceIdentifierType)
	if err != nil {
		return Empty{}, nil
	}
	getDeviceDetailsRequest, _ := getDeviceRequest.(GetDeviceDetailsRequest)
	return ListDeviceAttributesRequest(getDeviceDetailsRequest), nil
}

// nolint:bodyclose
func (r ListDeviceAttributesRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListDeviceAttributesResponse{res: res}, nil
}

func (r ListDeviceAttributesRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type SetDeviceAttributeRequest struct {
	req     *http.Request
	expects int
}

// SetDeviceAttribute builds a request to set an Attribute key to a certain value for a Device
func (c *Client) SetDeviceAttribute(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, attributeKey, attributeValue string) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL := makeURL(c.appEngineURL, "/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType))
	attributeMap := map[string]map[string]string{"attributes": {attributeKey: attributeValue}}
	payload, _ := makeBody(attributeMap)
	req := c.makeHTTPrequestWithContentType(http.MethodPatch, callURL, payload, "application/merge-patch+json")

	return SetDeviceAttributeRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r SetDeviceAttributeRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r SetDeviceAttributeRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type DeleteDeviceAttributeRequest struct {
	req     *http.Request
	expects int
}

// DeleteDeviceAttribute builds a request to delete an Attribute key and its value from a Device
func (c *Client) DeleteDeviceAttribute(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, attributeKey string) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL := makeURL(c.appEngineURL, "/v1/%s/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType))
	// We're using map[string]interface{} rather than map[string]string since we want to have null
	// rather than an empty string in the JSON payload, and this is the only way.
	attributeMap := map[string]map[string]interface{}{"attributes": {attributeKey: nil}}
	payload, _ := makeBody(attributeMap)
	req := c.makeHTTPrequestWithContentType(http.MethodPatch, callURL, payload, "application/merge-patch+json")

	return DeleteDeviceAttributeRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r DeleteDeviceAttributeRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r DeleteDeviceAttributeRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
