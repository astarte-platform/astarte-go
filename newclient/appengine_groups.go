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
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/astarte-platform/astarte-go/misc"
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
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups", realm))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)

	return ListGroupsRequest{req: req, expects: 200}, nil
}

func (r ListGroupsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return ListGroupsResponse{res: res}, nil
}

func (r ListGroupsRequest) ToCurl(c *Client) string {
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
		if !misc.IsValidAstarteDeviceID(deviceID) {
			return Empty{}, ErrInvalidDeviceID(deviceID)
		}
	}

	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups", realm))

	payload, _ := makeBody(DevicesAndGroup{GroupName: groupName, Devices: deviceIDList})
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload, c.token)
	return CreateGroupRequest{req: req, expects: 201}, nil
}

func (r CreateGroupRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return CreateGroupResponse{res: res}, nil
}

func (r CreateGroupRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

// ListGroupDevices builds a paginator to request a list of the devices that belong to a group.
func (c *Client) ListGroupDevices(realm, groupName string, pageSize int, format DeviceResultFormat) (Paginator, error) {
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups/%s/devices", realm, url.PathEscape(groupName)))

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
	if !misc.IsValidAstarteDeviceID(deviceID) {
		return Empty{}, ErrInvalidDeviceID(deviceID)
	}

	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups/%s/devices", realm, url.PathEscape(groupName)))

	payload, _ := makeBody(deviceIDPayload{Device: deviceID})
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload, c.token)

	return AddDeviceToGroupRequest{req: req, expects: 201}, nil
}

func (r AddDeviceToGroupRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r AddDeviceToGroupRequest) ToCurl(c *Client) string {
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
	if !misc.IsValidAstarteDeviceID(deviceID) {
		return Empty{}, ErrInvalidDeviceID(deviceID)
	}

	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/groups/%s/devices/%s", realm, url.PathEscape(groupName), deviceID))

	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil, c.token)

	return RemoveDeviceFromGroupRequest{req: req, expects: 204}, nil
}

func (r RemoveDeviceFromGroupRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {

		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r RemoveDeviceFromGroupRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
