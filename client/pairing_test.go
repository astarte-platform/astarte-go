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
	"testing"
)

func TestRegisterDevice(t *testing.T) {
	c, _ := getTestContext(t)
	registerDeviceCall, err := c.RegisterDevice(testRealmName, testDeviceID)
	if err != nil {
		t.Error(err)
	}
	res, err := registerDeviceCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	data, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	if data != testCredentialsSecret {
		t.Errorf("Failed credentials secret: %s\n", data)
	}
}

func TestUnregisterDevice(t *testing.T) {
	c, _ := getTestContext(t)
	unregisterDeviceCall, err := c.UnregisterDevice(testRealmName, testDeviceID)
	if err != nil {
		t.Error(err)
	}
	_, err = unregisterDeviceCall.Run(c)
	if err != nil {
		t.Error(err)
	}
}

func TestObtainNewMQTTv1CertificateForDevice(t *testing.T) {
	c, _ := getTestContext(t)
	getCertificateCall, _ := c.ObtainNewMQTTv1CertificateForDevice(testRealmName, testDeviceID, "a csr")
	getCertificateResponse, _ := getCertificateCall.Run(c)
	data, err := getCertificateResponse.Parse()
	if err != nil {
		t.Error(err)
	}
	if data != testClientCrt {
		t.Errorf("Failed certificate: %s\n", data)
	}
}

func TestGetMQTTv1ProtocolInformationForDevice(t *testing.T) {
	c, _ := getTestContext(t)
	getInfoCall, _ := c.GetMQTTv1ProtocolInformationForDevice(testRealmName, testDeviceID)
	getInfoResponse, _ := getInfoCall.Run(c)
	rawData, err := getInfoResponse.Parse()
	if err != nil {
		t.Error(err)
	}
	data, _ := rawData.(AstarteMQTTv1ProtocolInformation)
	if data.BrokerURL != testBrokerUrl {
		t.Errorf("Failed broker url: %s\n", data)
	}
}
