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
type listRealmsResponse struct {
	Res *http.Response
}

type getRealmResponse struct {
	Res *http.Response
}

type createRealmResponse struct {
	Res *http.Response
}

// Realm management
type listInterfacesResponse struct {
	Res *http.Response
}

type listInterfaceMajorVersionsResponse struct {
	Res *http.Response
}

type getInterfaceResponse struct {
	Res *http.Response
}

type installInterfaceResponse struct {
	Res *http.Response
}

type deleteInterfaceResponse struct {
	Res *http.Response
}
type updateInterfaceResponse struct {
	Res *http.Response
}

type listTriggersResponse struct {
	Res *http.Response
}

type getTriggerResponse struct {
	Res *http.Response
}

type installTriggerResponse struct {
	Res *http.Response
}

type deleteTriggerResponse struct {
	Res *http.Response
}

type GetNextDeviceListPageResponse struct {
	Res       *http.Response
	paginator *Paginator
}
