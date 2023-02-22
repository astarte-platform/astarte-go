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

	"github.com/astarte-platform/astarte-go/deviceid"
)

// resolveDeviceIdentifierType maps a deviceIdentifier and DeviceIdentifierType to a resolved
// DeviceIdentifierType (i.e. AstarteDeviceID or AstarteDeviceAlias). AutodiscoverDeviceIdentifier
// is resolved by checking if it's a valid Device ID, otherwise it's considered a Device Alias.
// AstarteDeviceAlias and AstarteDeviceID are returned as is.
func resolveDeviceIdentifierType(deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) DeviceIdentifierType {
	switch deviceIdentifierType {
	case AutodiscoverDeviceIdentifier:
		if deviceid.IsValid(deviceIdentifier) {
			return AstarteDeviceID
		}
		return AstarteDeviceAlias
	default:
		return deviceIdentifierType
	}
}

// devicePath accepts a deviceIdentifier and a resolved DeviceIdentifierType (i.e. AstarteDeviceID
// or AstarteDeviceAlias) and returns the path for that device. AutodiscoverDeviceIdentifier has to
// be resolved with resolveDeviceIdentifierType first
func devicePath(deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) string {
	switch deviceIdentifierType {
	case AstarteDeviceID:
		return fmt.Sprintf("devices/%v", deviceIdentifier)
	case AstarteDeviceAlias:
		return fmt.Sprintf("devices-by-alias/%v", deviceIdentifier)
	}
	return ""
}
