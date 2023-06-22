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
	"net/http"

	"github.com/astarte-platform/astarte-go/interfaces"
)

type AstarteResponse interface {
	// Parse reads the AstarteResponse returned by Run and returns either a well-typed
	// response payload or an error.
	Parse() (any, error)
	// Raw allows to supply a custom http Response handling function for the Astarte
	// response. The function does not need to close the response body.
	// Raw simply returns the value returned by the handling function.
	Raw(func(*http.Response) any) any
}

func (e Empty) Parse() (any, error)              { return nil, nil }
func (e Empty) Raw(func(*http.Response) any) any { return nil }

// Pairing

type RegisterDeviceResponse struct {
	res *http.Response
}

type NewDeviceCertificateResponse struct {
	res *http.Response
}

type Mqttv1DeviceInformationResponse struct {
	res *http.Response
}

// Housekeeping

type ListRealmsResponse struct {
	res *http.Response
}

type GetRealmResponse struct {
	res *http.Response
}

type CreateRealmResponse struct {
	res *http.Response
}

// Realm Management

type ListInterfacesResponse struct {
	res *http.Response
}

type ListInterfaceMajorVersionsResponse struct {
	res *http.Response
}

type GetInterfaceResponse struct {
	res *http.Response
}

type InstallInterfaceResponse struct {
	res *http.Response
}

type ListTriggersResponse struct {
	res *http.Response
}

type GetTriggerResponse struct {
	res *http.Response
}

type InstallTriggerResponse struct {
	res *http.Response
}

type ListTriggerDeliveryPoliciesResponse struct {
	res *http.Response
}

type GetTriggerDeliveryPolicyResponse struct {
	res *http.Response
}

type InstallTriggerDeliveryPolicyResponse struct {
	res *http.Response
}

// AppEngine

type GetNextDeviceListPageResponse struct {
	res       *http.Response
	paginator *Paginator
}

type GetDeviceIDFromAliasResponse struct {
	res *http.Response
}

type GetDeviceDetailsResponse struct {
	res *http.Response
}

type GetDeviceStatsResponse struct {
	res *http.Response
}

type ListDeviceInterfacesResponse struct {
	res *http.Response
}

type ListDeviceAliasesResponse struct {
	res *http.Response
}

type AddDeviceAliasResponse struct {
	res *http.Response
}

type ListDeviceAttributesResponse struct {
	res *http.Response
}

type GetNextDatastreamPageResponse struct {
	res       *http.Response
	paginator *Paginator
}

type GetDatastreamSnapshotResponse struct {
	res         *http.Response
	aggregation interfaces.AstarteInterfaceAggregation
}

type GetPropertiesResponse struct {
	res *http.Response
}

type ListGroupsResponse struct {
	res *http.Response
}

type CreateGroupResponse struct {
	res *http.Response
}

// General

type NoDataResponse struct {
	res *http.Response
}

// Parses data obtained by performing a request to Astarte which does not return data.
// The returned values do not matter.
func (r NoDataResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	return "", nil
}

func (r NoDataResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}
