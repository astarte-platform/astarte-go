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
	"net/http"
)

type AstarteResponse interface {
	// Parse reads the AstarteResponse returned by Run and returns either a well-typed
	// response payload or an error.
	Parse() (any, error)
	Raw()
}

// Used to represent an errror or an empty response.
type Empty struct{}

func (e Empty) Parse() (any, error) { return nil, nil }
func (e Empty) Raw()                {}

// Pairing
type registerDeviceResponse struct {
	Res *http.Response
}

type unregisterDeviceResponse struct {
	Res *http.Response
}

type newDeviceCertificateResponse struct {
	Res *http.Response
}

type mqttv1DeviceInformationResponse struct {
	Res *http.Response
}

// Housekeeping
type ListRealmsResponse struct {
	Res *http.Response
}

type GetRealmResponse struct {
	Res *http.Response
}

type CreateRealmResponse struct {
	Res *http.Response
}

// Realm management
type ListInterfacesResponse struct {
	Res *http.Response
}

type ListInterfaceMajorVersionsResponse struct {
	Res *http.Response
}

type GetInterfaceResponse struct {
	Res *http.Response
}

type InstallInterfaceResponse struct {
	Res *http.Response
}

type DeleteInterfaceResponse struct {
	Res *http.Response
}
type UpdateInterfaceResponse struct {
	Res *http.Response
}

type ListTriggersResponse struct {
	Res *http.Response
}

type GetTriggerResponse struct {
	Res *http.Response
}

type InstallTriggerResponse struct {
	Res *http.Response
}

type DeleteTriggerResponse struct {
	Res *http.Response
}

type GetNextDeviceListPageResponse struct {
	Res       *http.Response
	paginator *Paginator
}
