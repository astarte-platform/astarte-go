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
	"io"

	"github.com/tidwall/gjson"
)

// Parses data obtained by performing a request to register a device.
// Returns the new credentials secret as a string.
func (r registerDeviceResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	value := gjson.GetBytes(b, "data.credentials_secret").String()
	return value, nil
}
func (e registerDeviceResponse) Raw() {}

// Parses data obtained by performing a request to unregister a device.
// The returned values do not matter.
func (r unregisterDeviceResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	return "", nil
}
func (e unregisterDeviceResponse) Raw() {}

// Parses data obtained by performing a request to for a new device certificate.
// Returns the new device certificate as a PEM-encoded string.
func (r newDeviceCertificateResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	value := gjson.GetBytes(b, "data.client_crt").String()
	return value, nil
}
func (e newDeviceCertificateResponse) Raw() {}

// Parses data obtained by performing a request for connection information
// for a newly registered device.
// Returns the Astarte broker URL as a string.
func (r mqttv1DeviceInformationResponse) Parse() (any, error) {
	defer r.Res.Body.Close()
	b, _ := io.ReadAll(r.Res.Body)
	value := gjson.GetBytes(b, "data.broker_url").String()
	return value, nil
}
func (e mqttv1DeviceInformationResponse) Raw() {}
