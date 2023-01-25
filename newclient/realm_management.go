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

package newclient

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/astarte-platform/astarte-go/interfaces"
	"moul.io/http2curl"
)

type listInterfacesRequest struct {
	req     *http.Request
	expects int
}

// ListInterfaces builds a request to return all interfaces in a Realm.
func (c *Client) ListInterfaces(realm string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces", realm))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return listInterfacesRequest{req: req, expects: 200}, nil
}

func (r listInterfacesRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return listInterfacesResponse{Res: res}, nil
}

func (r listInterfacesRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type listInterfaceMajorVersionsRequest struct {
	req     *http.Request
	expects int
}

// ListInterfaceMajorVersions builds a request to return all available major versions for a given Interface in a Realm.
func (c *Client) ListInterfaceMajorVersions(realm string, interfaceName string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s", realm, interfaceName))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return listInterfaceMajorVersionsRequest{req: req, expects: 200}, nil
}

func (r listInterfaceMajorVersionsRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return listInterfaceMajorVersionsResponse{Res: res}, nil
}

func (r listInterfaceMajorVersionsRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type getInterfaceRequest struct {
	req     *http.Request
	expects int
}

// GetInterface builds a request retrieve an interface, identified by a Major version, in a Realm.
func (c *Client) GetInterface(realm string, interfaceName string, interfaceMajor int) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s/%v", realm, interfaceName, interfaceMajor))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return getInterfaceRequest{req: req, expects: 200}, nil
}

func (r getInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return getInterfaceResponse{Res: res}, nil
}

func (r getInterfaceRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type installInterfaceRequest struct {
	req     *http.Request
	expects int
}

// InstallInterface builds a request to install a new major version of an Interface into the Realm.
func (c *Client) InstallInterface(realm string, interfacePayload interfaces.AstarteInterface) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces", realm))

	payload, _ := makeBody(interfacePayload)
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload, c.token)
	return installInterfaceRequest{req: req, expects: 201}, nil
}

func (r installInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return installInterfaceResponse{Res: res}, nil
}

func (r installInterfaceRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type deleteInterfaceRequest struct {
	req     *http.Request
	expects int
}

// DeleteInterface builds a request to delete a major version of an Interface into the Realm.
func (c *Client) DeleteInterface(realm string, interfaceName string, interfaceMajor int) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s/%v", realm, interfaceName, interfaceMajor))

	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil, c.token)
	return deleteInterfaceRequest{req: req, expects: 204}, nil
}

func (r deleteInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return deleteInterfaceResponse{Res: res}, nil
}

func (r deleteInterfaceRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type updateInterfaceRequest struct {
	req     *http.Request
	expects int
}

// UpdateInterface builds a request to update an existing major version of an Interface to a new minor.
func (c *Client) UpdateInterface(realm string, interfaceName string, interfaceMajor int, interfacePayload interfaces.AstarteInterface) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/interfaces/%s/%v", realm, interfaceName, interfaceMajor))

	payload, _ := makeBody(interfacePayload)
	req := c.makeHTTPrequest(http.MethodPut, callURL, payload, c.token)
	return updateInterfaceRequest{req: req, expects: 204}, nil
}

func (r updateInterfaceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return updateInterfaceResponse{Res: res}, nil
}

func (r updateInterfaceRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type listTriggersRequest struct {
	req     *http.Request
	expects int
}

// ListTriggers builds a request to return all triggers in a Realm.
func (c *Client) ListTriggers(realm string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers", realm))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return listTriggersRequest{req: req, expects: 200}, nil
}

func (r listTriggersRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return listTriggersResponse{Res: res}, nil
}

func (r listTriggersRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type getTriggerRequest struct {
	req     *http.Request
	expects int
}

// GetTrigger builds a request to return a trigger installed in a Realm.
func (c *Client) GetTrigger(realm string, triggerName string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers/%s", realm, triggerName))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)
	return getTriggerRequest{req: req, expects: 200}, nil
}

func (r getTriggerRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return getTriggerResponse{Res: res}, nil
}

func (r getTriggerRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type installTriggerRequest struct {
	req     *http.Request
	expects int
}

// InstallTrigger builds a request to install a Trigger into the Realm.
func (c *Client) InstallTrigger(realm string, triggerPayload any) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers", realm))

	payload, _ := makeBody(triggerPayload)
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload, c.token)
	return installTriggerRequest{req: req, expects: 201}, nil
}

func (r installTriggerRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return installTriggerResponse{Res: res}, nil
}

func (r installTriggerRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type deleteTriggerRequest struct {
	req     *http.Request
	expects int
}

// DeleteTrigger builds a request to delete a Trigger from the Realm.
func (c *Client) DeleteTrigger(realm string, triggerName string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.realmManagementURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/triggers/%s", realm, triggerName))

	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil, c.token)
	return deleteTriggerRequest{req: req, expects: 204}, nil
}

func (r deleteTriggerRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return empty{}, err
	}
	if res.StatusCode != r.expects {
		return empty{}, ErrDifferentStatusCode
	}
	return deleteTriggerResponse{Res: res}, nil
}

func (r deleteTriggerRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
