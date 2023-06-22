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

	"github.com/astarte-platform/astarte-go/deviceid"
	"moul.io/http2curl"
)

// DevicesAndGroup maps to the JSON object returned by a Create Group call to AppEngine API.
type DevicesAndGroup struct {
	GroupName string   `json:"group_name"`
	Devices   []string `json:"devices"`
}

type deviceIDPayload struct {
	Device string `json:"device_id"`
}

type ListGroupsRequest struct {
	req     *http.Request
	expects int
}

// ListGroups builds a request to list the groups in a Realm.
func (c *Client) ListGroups(realm string) (AstarteRequest, error) {
	callURL := makeURL(c.appEngineURL, "/v1/%s/groups", realm)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return ListGroupsRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r ListGroupsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListGroupsResponse{res: res}, nil
}

func (r ListGroupsRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type CreateGroupRequest struct {
	req     *http.Request
	expects int
}

// CreateGroup builds a request to create a group with the given deviceIDList in the Realm.
// Only valid Astarte device IDs can be used when adding devices to a group.
func (c *Client) CreateGroup(realm, groupName string, deviceIDList []string) (AstarteRequest, error) {
	for _, deviceID := range deviceIDList {
		if !deviceid.IsValid(deviceID) {
			return Empty{}, ErrInvalidDeviceID(deviceID)
		}
	}

	callURL := makeURL(c.appEngineURL, "/v1/%s/groups", realm)
	payload, _ := makeBody(DevicesAndGroup{GroupName: groupName, Devices: deviceIDList})
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload)

	return CreateGroupRequest{req: req, expects: 201}, nil
}

// nolint:bodyclose
func (r CreateGroupRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return CreateGroupResponse{res: res}, nil
}

func (r CreateGroupRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

// ListGroupDevices builds a paginator to request a list of the devices that belong to a group.
func (c *Client) ListGroupDevices(realm, groupName string, pageSize int, format DeviceResultFormat) (Paginator, error) {
	callURL := makeURL(c.appEngineURL, "/v1/%s/groups/%s/devices", realm, url.PathEscape(groupName))
	paginator, err := c.GetDeviceListPaginator(realm, pageSize, format)
	if err != nil {
		return &DeviceListPaginator{}, err
	}

	deviceListPaginator := paginator.(*DeviceListPaginator)
	deviceListPaginator.baseURL = callURL

	return deviceListPaginator, nil
}

type AddDeviceToGroupRequest struct {
	req     *http.Request
	expects int
}

// AddDeviceToGroup builds a request to add a device to a group.
// Only valid Astarte device IDs can be used when adding a device to a group.
func (c *Client) AddDeviceToGroup(realm, groupName, deviceID string) (AstarteRequest, error) {
	if !deviceid.IsValid(deviceID) {
		return Empty{}, ErrInvalidDeviceID(deviceID)
	}

	callURL := makeURL(c.appEngineURL, "/v1/%s/groups/%s/devices", realm, url.PathEscape(groupName))
	payload, _ := makeBody(deviceIDPayload{Device: deviceID})
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload)

	return AddDeviceToGroupRequest{req: req, expects: 201}, nil
}

// nolint:bodyclose
func (r AddDeviceToGroupRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r AddDeviceToGroupRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type RemoveDeviceFromGroupRequest struct {
	req     *http.Request
	expects int
}

// RemoveDeviceFromGroup builds a request to removes a device from the group.
// Only valid Astarte device IDs can be used when removing a device from a group.
func (c *Client) RemoveDeviceFromGroup(realm, groupName, deviceID string) (AstarteRequest, error) {
	if !deviceid.IsValid(deviceID) {
		return Empty{}, ErrInvalidDeviceID(deviceID)
	}

	callURL := makeURL(c.appEngineURL, "/v1/%s/groups/%s/devices/%s", realm, url.PathEscape(groupName), deviceID)
	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil)

	return RemoveDeviceFromGroupRequest{req: req, expects: 204}, nil
}

// nolint:bodyclose
func (r RemoveDeviceFromGroupRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode(r.expects, res.StatusCode)
	}
	return NoDataResponse{res: res}, nil
}

func (r RemoveDeviceFromGroupRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
