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
	"net/http"
	"strconv"

	"github.com/astarte-platform/astarte-go/interfaces"
	"moul.io/http2curl"
)

type ListInterfacesRequest struct {
	req     *http.Request
	expects int
}

// ListInterfaces builds a request to return all interfaces in a Realm.
func (c *Client) ListInterfaces(realm string) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/interfaces", realm)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return ListInterfacesRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r ListInterfacesRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListInterfacesResponse{res: res}, nil
}

func (r ListInterfacesRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type ListInterfaceMajorVersionsRequest struct {
	req     *http.Request
	expects int
}

// ListInterfaceMajorVersions builds a request to return all available major versions for a given Interface in a Realm.
func (c *Client) ListInterfaceMajorVersions(realm string, interfaceName string) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/interfaces/%s", realm, interfaceName)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return ListInterfaceMajorVersionsRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r ListInterfaceMajorVersionsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListInterfaceMajorVersionsResponse{res: res}, nil
}

func (r ListInterfaceMajorVersionsRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type GetInterfaceRequest struct {
	req     *http.Request
	expects int
}

// GetInterface builds a request retrieve an interface, identified by a Major version, in a Realm.
func (c *Client) GetInterface(realm string, interfaceName string, interfaceMajor int) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/interfaces/%s/%s", realm, interfaceName, fmt.Sprintf("%v", interfaceMajor))
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return GetInterfaceRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r GetInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return GetInterfaceResponse{res: res}, nil
}

func (r GetInterfaceRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type InstallInterfaceRequest struct {
	req     *http.Request
	expects int
}

// InstallInterface builds a request to install a new major version of an Interface into the Realm.
func (c *Client) InstallInterface(realm string, interfacePayload interfaces.AstarteInterface, isAsync bool) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/interfaces", realm)

	if !isAsync {
		query := map[string]string{"async_operation": strconv.FormatBool(false)}
		callURL = setupURLQuery(callURL, query)
	}

	payload, _ := makeBody(interfacePayload)
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload)

	return InstallInterfaceRequest{req: req, expects: 201}, nil
}

// nolint:bodyclose
func (r InstallInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return InstallInterfaceResponse{res: res}, nil
}

func (r InstallInterfaceRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type DeleteInterfaceRequest struct {
	req     *http.Request
	expects int
}

// DeleteInterface builds a request to delete a major version of an Interface into the Realm.
func (c *Client) DeleteInterface(realm string, interfaceName string, interfaceMajor int) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/interfaces/%s/%s", realm, interfaceName, fmt.Sprintf("%v", interfaceMajor))
	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil)

	return DeleteInterfaceRequest{req: req, expects: 204}, nil
}

// nolint:bodyclose
func (r DeleteInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r DeleteInterfaceRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type UpdateInterfaceRequest struct {
	req     *http.Request
	expects int
}

// UpdateInterface builds a request to update an existing major version of an Interface to a new minor.
func (c *Client) UpdateInterface(realm string, interfaceName string, interfaceMajor int, interfacePayload interfaces.AstarteInterface, isAsync bool) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/interfaces/%s/%s", realm, interfaceName, fmt.Sprintf("%v", interfaceMajor))

	if !isAsync {
		query := map[string]string{"async_operation": strconv.FormatBool(false)}
		callURL = setupURLQuery(callURL, query)
	}

	payload, _ := makeBody(interfacePayload)
	req := c.makeHTTPrequest(http.MethodPut, callURL, payload)

	return UpdateInterfaceRequest{req: req, expects: 204}, nil
}

// nolint:bodyclose
func (r UpdateInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r UpdateInterfaceRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type ListTriggersRequest struct {
	req     *http.Request
	expects int
}

// ListTriggers builds a request to return all triggers in a Realm.
func (c *Client) ListTriggers(realm string) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/triggers", realm)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return ListTriggersRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r ListTriggersRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return ListTriggersResponse{res: res}, nil
}

func (r ListTriggersRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type GetTriggerRequest struct {
	req     *http.Request
	expects int
}

// GetTrigger builds a request to return a trigger installed in a Realm.
func (c *Client) GetTrigger(realm string, triggerName string) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/triggers/%s", realm, triggerName)
	req := c.makeHTTPrequest(http.MethodGet, callURL, nil)

	return GetTriggerRequest{req: req, expects: 200}, nil
}

// nolint:bodyclose
func (r GetTriggerRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return GetTriggerResponse{res: res}, nil
}

func (r GetTriggerRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type InstallTriggerRequest struct {
	req     *http.Request
	expects int
}

// InstallTrigger builds a request to install a Trigger into the Realm.
func (c *Client) InstallTrigger(realm string, triggerPayload any) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/triggers", realm)
	payload, _ := makeBody(triggerPayload)
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload)

	return InstallTriggerRequest{req: req, expects: 201}, nil
}

// nolint:bodyclose
func (r InstallTriggerRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return InstallTriggerResponse{res: res}, nil
}

func (r InstallTriggerRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type DeleteTriggerRequest struct {
	req     *http.Request
	expects int
}

// DeleteTrigger builds a request to delete a Trigger from the Realm.
func (c *Client) DeleteTrigger(realm string, triggerName string) (AstarteRequest, error) {
	callURL := makeURL(c.realmManagementURL, "/v1/%s/triggers/%s", realm, triggerName)
	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil)

	return DeleteTriggerRequest{req: req, expects: 204}, nil
}

// nolint:bodyclose
func (r DeleteTriggerRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return runAstarteRequestError(res, r.expects)
	}
	return NoDataResponse{res: res}, nil
}

func (r DeleteTriggerRequest) ToCurl(_ *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
