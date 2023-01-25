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
	"fmt"
	"net/http"
	"net/url"
	"path"

	"moul.io/http2curl"
)

type registerDevicePayload struct {
	HwID string `json:"hw_id"`
}

type getMQTTv1CertificatePayload struct {
	CSR string `json:"csr"`
}

type RegisterDeviceRequest struct {
	req     *http.Request
	expects int
}

// RegisterDevice builds a request to register a new device into the Realm.
// TODO: add support for initial_introspection
func (c *Client) RegisterDevice(realm string, deviceID string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/agent/devices", realm))

	// TODO check err
	payload, _ := makeBody(registerDevicePayload{HwID: deviceID})
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload, c.token)
	return RegisterDeviceRequest{req: req, expects: 201}, nil
}

func (r RegisterDeviceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return RegisterDeviceResponse{res: res}, nil
}

func (r RegisterDeviceRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type UnregisterDeviceRequest struct {
	req     *http.Request
	expects int
}

// UnregisterDevice builds a request to reset the registration state of a device.
// Once the request is run, this makes it possible to register it again.
// All data belonging to the device will be left as is in Astarte.
func (c *Client) UnregisterDevice(realm string, deviceID string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/agent/devices/%s", realm, deviceID))

	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil, c.token)
	return UnregisterDeviceRequest{req: req, expects: 204}, nil
}

func (r UnregisterDeviceRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r UnregisterDeviceRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type NewDeviceCertificateRequest struct {
	req     *http.Request
	expects int
}

// ObtainNewMQTTv1CertificateForDevice builds a request for retrieving a valid SSL Certificate for Devices
// running on astarte_mqtt_v1.
// This API is meant to be called by the device, and the Client that executes (Runs) the request needs to
// have the Device's Credentials Secret as its token.
func (c *Client) ObtainNewMQTTv1CertificateForDevice(realm, deviceID, csr string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s/protocols/astarte_mqtt_v1/credentials", realm, deviceID))

	payload, _ := makeBody(getMQTTv1CertificatePayload{CSR: csr})
	req := c.makeHTTPrequest(http.MethodPost, callURL, payload, c.token)

	return NewDeviceCertificateRequest{req: req, expects: 201}, nil
}

func (r NewDeviceCertificateRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r NewDeviceCertificateRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type Mqttv1DeviceInformationRequest struct {
	req     *http.Request
	expects int
}

// GetMQTTv1ProtocolInformationForDevice builds a request for retrieving protocol information (such as
// the broker URL) for devices running on astarte_mqtt_v1.
// This API is meant to be called by the device, and the Client that executes (Runs) the request needs to
// have the Device's Credentials Secret as its token.
func (c *Client) GetMQTTv1ProtocolInformationForDevice(realm, deviceID string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s", realm, deviceID))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)

	return Mqttv1DeviceInformationRequest{req: req, expects: 200}, nil
}

func (r Mqttv1DeviceInformationRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return Mqttv1DeviceInformationResponse{res: res}, nil
}

func (r Mqttv1DeviceInformationRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
