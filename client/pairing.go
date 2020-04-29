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

// PairingService is the API Client for Pairing API
type PairingService struct {
	client     *Client
	pairingURL *url.URL
}

// RegisterDevice registers a new device into the Realm.
// Returns the Credential Secret of the Device when successful.
// TODO: add support for initial_introspection
func (s *PairingService) RegisterDevice(realm string, deviceID string) (string, error) {
	callURL, _ := url.Parse(s.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/agent/devices", realm))

	var requestBody struct {
		HwID string `json:"hw_id"`
	}
	requestBody.HwID = deviceID

	ret := deviceRegistrationResponse{}
	err := s.client.genericJSONDataAPIPostWithResponse(&ret, callURL.String(), requestBody, 201)

	return ret.CredentialsSecret, err
}

// UnregisterDevice resets the registration state of a device. This makes it possible to register it again.
// All data belonging to the device will be left as is in Astarte.
func (s *PairingService) UnregisterDevice(realm string, deviceID string) error {
	callURL, _ := url.Parse(s.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/agent/devices/%s", realm, deviceID))

	err := s.client.genericJSONDataAPIDelete(callURL.String(), 204)
	if err != nil {
		return err
	}

	return nil
}

// ObtainNewMQTTv1CertificateForDevice returns a valid SSL Certificate for Devices running on astarte_mqtt_v1.
// This API is meant to be called by the device, and your Client needs to have the Device's Credentials Secret
// as its token. Always call SetToken with the Credentials Secret before calling this function.
func (s *PairingService) ObtainNewMQTTv1CertificateForDevice(realm, deviceID, csr string) (string, error) {
	callURL, _ := url.Parse(s.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s/protocols/astarte_mqtt_v1/credentials", realm, deviceID))

	var requestBody struct {
		CSR string `json:"csr"`
	}
	requestBody.CSR = csr

	ret := getMQTTv1CertificateResponse{}
	err := s.client.genericJSONDataAPIPostWithResponse(&ret, callURL.String(), requestBody, 201)

	return ret.ClientCertificate, err
}

// GetMQTTv1ProtocolInformationForDevice returns protocol information (such as the broker URL) for devices running
// on astarte_mqtt_v1.
// This API is meant to be called by the device, and your Client needs to have the Device's Credentials Secret
// as its token. Always call SetToken with the Credentials Secret before calling this function.
func (s *PairingService) GetMQTTv1ProtocolInformationForDevice(realm, deviceID string) (AstarteMQTTv1ProtocolInformation, error) {
	callURL, _ := url.Parse(s.pairingURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/devices/%s", realm, deviceID))

	ret := AstarteMQTTv1ProtocolInformation{}
	err := s.client.genericJSONDataAPIGET(&ret, callURL.String(), 200)

	return ret, err
}
