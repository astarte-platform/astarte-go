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
	"encoding/json"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

type AstarteMQTTv1ProtocolInformation struct {
	BrokerURL string `json:"broker_url"`
}

// Parses data obtained by performing a request to register a device.
// Returns the new credentials secret as a string.
func (r RegisterDeviceResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	value := gjson.GetBytes(b, "data.credentials_secret").String()
	return value, nil
}
func (r RegisterDeviceResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request to for a new device certificate.
// Returns the new device certificate as a PEM-encoded string.
func (r NewDeviceCertificateResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	value := gjson.GetBytes(b, "data.client_crt").String()
	return value, nil
}
func (r NewDeviceCertificateResponse) Raw() *http.Response {
	return r.res
}

// Parses data obtained by performing a request for connection information
// for a newly registered device.
// Returns the information as an AstarteMQTTv1ProtocolInformation struct.
func (r Mqttv1DeviceInformationResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data").Raw
	value := AstarteMQTTv1ProtocolInformation{}
	_ = json.Unmarshal([]byte(data), &value)
	return value, nil
}
func (r Mqttv1DeviceInformationResponse) Raw() *http.Response {
	return r.res
}
