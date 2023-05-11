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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astarte-platform/astarte-go/interfaces"
	"github.com/iancoleman/orderedmap"
	"github.com/nqd/flat"
	"github.com/tidwall/gjson"
)

type Paginator interface {
	GetNextPage() (AstarteRequest, error)
	GetPageSize() int
	HasNextPage() bool
	Rewind()

	computePageState(rawData []byte)
	parseData(rawData []byte) any
}

// DeviceResultFormat represents the format of the Device returned in the Device list.
type DeviceResultFormat int

const (
	// DeviceIDFormat means the Paginator will return a list of strings
	// representing the Device ID of the Devices.
	DeviceIDFormat DeviceResultFormat = iota
	// DeviceDetailsFormat means the Paginator will return a list of
	// DeviceDetails structs
	DeviceDetailsFormat
)

// Links is a struct that represent the links metadata returned by Astarte API.
// This metadata is used in Astarte APIs to perform pagination, allowing the
// user to simply follow the Next link, if any, to fetch the next page.
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
}

// DeviceInterfaceIntrospection represents a single entry in a Device Introspection array retrieved
// from DeviceDetails.
type DeviceInterfaceIntrospection struct {
	Name              string `json:"name,omitempty"`
	Major             int    `json:"major"`
	Minor             int    `json:"minor"`
	ExchangedMessages uint64 `json:"exchanged_msgs,omitempty"`
	ExchangedBytes    uint64 `json:"exchanged_bytes,omitempty"`
}

// DeviceDetails maps to the JSON object returned by a Device Details call to AppEngine API.
type DeviceDetails struct {
	TotalReceivedMessages    int64                                   `json:"total_received_msgs"`
	TotalReceivedBytes       uint64                                  `json:"total_received_bytes"`
	LastSeenIP               net.IP                                  `json:"last_seen_ip"`
	LastDisconnection        time.Time                               `json:"last_disconnection"`
	LastCredentialsRequestIP net.IP                                  `json:"last_credentials_request_ip"`
	LastConnection           time.Time                               `json:"last_connection"`
	DeviceID                 string                                  `json:"id"`
	FirstRegistration        time.Time                               `json:"first_registration"`
	FirstCredentialsRequest  time.Time                               `json:"first_credentials_request"`
	CredentialsInhibited     bool                                    `json:"credentials_inhibited"`
	Connected                bool                                    `json:"connected"`
	Introspection            map[string]DeviceInterfaceIntrospection `json:"introspection"`
	Aliases                  map[string]string                       `json:"aliases"`
	PreviousInterfaces       []DeviceInterfaceIntrospection          `json:"previous_interfaces,omitempty"`
	Attributes               map[string]string                       `json:"attributes,omitempty"`
}

// DevicesStats maps to the JSON object returned by a Device Stats call to AppEngine API.
type DevicesStats struct {
	TotalDevices     int64 `json:"total_devices"`
	ConnectedDevices int64 `json:"connected_devices"`
}

// Parses data obtained by performing a request a DeviceListPaginator page
// and sets up the paginator for retrieving the next page.
// Returns the page as an array of strings or DeviceDetails, depending on the format specified in the paginator.
func (r GetNextDeviceListPageResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)

	// Golang I hate you so much
	paginator := (*r.paginator).(*DeviceListPaginator)

	data := paginator.parseData(b)
	paginator.computePageState(b)

	return data, nil
}

func updatePaginator(p Paginator, res *http.Response) io.ReadCloser {
	b, _ := io.ReadAll(res.Body)

	// Create a helper ReadCloser to keep a copy of the response body
	readCloser := io.NopCloser(bytes.NewReader(b))

	// update the state of the paginator
	p.computePageState(b)

	// and return a copy of the body
	return readCloser
}

// Raw allows to supply a custom http Response handling function for the Astarte
// response. The handling function must not close the body of the response. Moreover,
// Raw sets up the paginator for retrieving the next page.
// Raw simply returns the value returned by the handling function.
func (r GetNextDeviceListPageResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()

	p := (*r.paginator).(*DeviceListPaginator)
	r.res.Body = updatePaginator(p, r.res)

	return f(r.res)
}

func (d *DeviceListPaginator) parseData(rawData []byte) any {
	data := gjson.GetBytes(rawData, "data").Array()
	switch d.format {
	case DeviceIDFormat:
		ret := []string{}
		for _, v := range data {
			ret = append(ret, v.Str)
		}
		return ret
	case DeviceDetailsFormat:
		ret := []DeviceDetails{}
		for _, v := range data {
			details := DeviceDetails{}
			_ = json.Unmarshal([]byte(v.Raw), &details)
			ret = append(ret, details)
		}
		return ret
	// we'll never get there as there are only 2 formats
	default:
		return nil
	}
}

func (d *DeviceListPaginator) computePageState(rawData []byte) {
	links := Links{}
	_ = json.Unmarshal(rawData, &links)
	if links.Next == "" {
		d.hasNextPage = false
	} else {
		d.hasNextPage = true
		parsedLinks, _ := url.Parse(links.Next)
		d.nextQuery = parsedLinks.Query()
	}
}

// Parses data obtained by performing a request a Device ID from alias.
// Returns the device ID as a string.
func (r GetDeviceIDFromAliasResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data")
	details := DeviceDetails{}
	_ = json.Unmarshal([]byte(data.Raw), &details)
	return details.DeviceID, nil
}

func (r GetDeviceIDFromAliasResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request device details.
// Returns details as a DeviceDetails structure.
func (r GetDeviceDetailsResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data")
	details := DeviceDetails{}
	_ = json.Unmarshal([]byte(data.Raw), &details)
	return details, nil
}

func (r GetDeviceDetailsResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request a device introspection.
// Returns the list of interface names as an array of strings.
func (r ListDeviceInterfacesResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data").Array()
	interfaces := []string{}
	for _, v := range data {
		interfaces = append(interfaces, v.Str)
	}
	return interfaces, nil
}

func (r ListDeviceInterfacesResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request device's aliases.
// Returns the list of aliases as an array of strings.
func (r ListDeviceAliasesResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data.aliases").Array()
	aliases := []string{}
	for _, v := range data {
		aliases = append(aliases, v.Str)
	}
	return aliases, nil
}

func (r ListDeviceAliasesResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request device's attributes.
// Returns the attributes as a map strings to strings.
func (r ListDeviceAttributesResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data.attributes").Map()
	attributes := map[string]string{}
	for k, v := range data {
		attributes[k] = v.Str
	}
	return attributes, nil
}

func (r ListDeviceAttributesResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request for devices stats.
// Returns the stats as a DevicesStats struct.
func (r GetDeviceStatsResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data")
	stats := DevicesStats{}
	_ = json.Unmarshal([]byte(data.Raw), &stats)
	return stats, nil
}

func (r GetDeviceStatsResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// DatastreamIndividualValue represent one single Datastream value on an interface with Individual aggregation.
type DatastreamIndividualValue struct {
	Value              interface{} `json:"value"`
	Timestamp          time.Time   `json:"timestamp"`
	ReceptionTimestamp time.Time   `json:"reception_timestamp,omitempty"`
}

// DatastreamIndividualValue represent one Datastream value on an interface with Object aggregation.
type DatastreamObjectValue struct {
	Values    orderedmap.OrderedMap
	Timestamp time.Time
}

// PropertyValue represent the Property value on a properties interface.
type PropertyValue any

// UnmarshalJSON unmarshals a quoted json string to a DatastreamObjectValue
func (s *DatastreamObjectValue) UnmarshalJSON(b []byte) error {
	var j orderedmap.OrderedMap
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	// just to check that JSON did not curse the timestamo
	timestampInterface, _ := j.Get("timestamp")
	switch v := timestampInterface.(type) {
	case time.Time:
		s.Timestamp = v
	case string:
		var err error
		s.Timestamp, err = time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return err
		}
	}

	j.Delete("timestamp")
	s.Values = j

	return nil
}

// Parses data obtained by performing a request for a DatastreamPaginator page
// and sets up the paginator for retrieving the next page.
// According to the interface's aggregation and path, the return value can be one of:
// []DatastreamObjectValue, map[string][]DatastreamObjectValue, []DatastreamIndividualValue,
// map[string]DatastreamIndividualValue.
func (r GetNextDatastreamPageResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)

	// Golang I hate you so much
	paginator := (*r.paginator).(*DatastreamPaginator)

	data := paginator.parseData(b)
	paginator.computePageState(b)

	return data, nil
}

// Raw allows to supply a custom http Response handling function for the Astarte
// response. The handling function must not close the body of the response. Moreover,
// Raw sets up the paginator for retrieving the next page.
// Raw simply returns the value returned by the handling function.
func (r GetNextDatastreamPageResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()

	p := (*r.paginator).(*DatastreamPaginator)
	r.res.Body = updatePaginator(p, r.res)

	return f(r.res)
}

func (d *DatastreamPaginator) parseData(rawData []byte) any {
	data := gjson.GetBytes(rawData, "data").Raw
	jsonData := gjson.ParseBytes([]byte(data))
	return parseDatastream(jsonData, d.aggregation)
}

func parseDatastream(jsonData gjson.Result, aggregation interfaces.AstarteInterfaceAggregation) any {
	// handle the case of individual aggregation
	if aggregation == interfaces.IndividualAggregation {
		return parseDatastreamWithIndividualAggregation(jsonData)
	}

	// handle object aggregation
	return parseDatastreamWithObjectAggregation(jsonData)
}

func parseDatastreamWithObjectAggregation(jsonData gjson.Result) any {
	if jsonData.IsArray() {
		objectValues := []DatastreamObjectValue{}
		data := jsonData.Array()
		for _, v := range data {
			value := DatastreamObjectValue{}
			_ = json.Unmarshal([]byte(v.Raw), &value)
			objectValues = append(objectValues, value)
		}
		return objectValues
	}
	// if not an array, it must be an object
	obj := jsonData.Value().(map[string]interface{})

	// now we need to flatten the object so that the common portion of the path can be factored out
	// from each mapping
	flattened, _ := flat.Flatten(obj, &flat.Options{Safe: true, Delimiter: "."})

	keys := []string{}
	for k := range flattened {
		components := strings.Split(k, ".")
		var theKey string
		if len(components) > 1 {
			theKey = strings.Join(components[:len(components)-1], ".")
		} else {
			theKey = k
		}
		keys = append(keys, theKey)
	}
	keys = removeDuplicateStr(keys)

	// and once we have all the keys, we can get the object values
	rawObjectValues := gjson.GetMany(jsonData.Raw, keys...)

	ret := map[string][]DatastreamObjectValue{}
	for i, item := range rawObjectValues {
		values := []DatastreamObjectValue{}
		value := DatastreamObjectValue{}

		k := fmt.Sprintf("/%s", strings.ReplaceAll(keys[i], ".", "/"))

		if item.IsArray() {
			_ = json.Unmarshal([]byte(item.Raw), &values)
			ret[k] = append(ret[k], values...)
		} else {
			_ = json.Unmarshal([]byte(item.Raw), &value)
			ret[k] = append(ret[k], value)
		}
	}
	return ret
}

func parseDatastreamWithIndividualAggregation(jsonData gjson.Result) any {
	// first, we check if the complete timeseries is returned
	individualValues := []DatastreamIndividualValue{}
	if jsonData.IsArray() {
		data := jsonData.Array()
		for _, v := range data {
			value := DatastreamIndividualValue{}
			_ = json.Unmarshal([]byte(v.Raw), &value)
			individualValues = append(individualValues, value)
		}
		return individualValues
	}

	// if it's not a timeseries, it must be a snapshot (objects are returned)
	obj := jsonData.Value().(map[string]interface{})

	// now we need to flatten the object so that the common portion of the path can be factored out
	// from each mapping
	flattened, _ := flat.Flatten(obj, &flat.Options{Safe: true, Delimiter: "."})

	keys := []string{}
	for k := range flattened {
		components := strings.Split(k, ".")
		theKey := strings.Join(components[:len(components)-1], ".")
		keys = append(keys, theKey)
	}
	keys = removeDuplicateStr(keys)

	// and once we have all the keys, we can get the object values
	rawIndividualValues := gjson.GetMany(jsonData.Raw, keys...)

	ret := map[string]DatastreamIndividualValue{}
	for i, item := range rawIndividualValues {
		value := DatastreamIndividualValue{}
		_ = json.Unmarshal([]byte(item.Raw), &value)
		k := fmt.Sprintf("/%s", strings.ReplaceAll(keys[i], ".", "/"))
		ret[k] = value
	}
	return ret
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := map[string]struct{}{}
	list := []string{}
	for _, item := range strSlice {
		if _, ok := allKeys[item]; !ok {
			allKeys[item] = struct{}{}
			list = append(list, item)
		}
	}
	return list
}

func (d *DatastreamPaginator) computePageState(rawData []byte) {
	data := gjson.GetBytes(rawData, "data").Array()
	resultLength := len(data)
	if resultLength < d.pageSize {
		d.hasNextPage = false
	} else {
		d.hasNextPage = true
		d.firstPage = false
		d.updateTimestampValues(data[resultLength-1])
	}
}

func (d *DatastreamPaginator) updateTimestampValues(updateValue gjson.Result) {
	if updateValue.Get("value").Exists() {
		val := DatastreamIndividualValue{}
		_ = json.Unmarshal([]byte(updateValue.Raw), &val)
		switch d.resultSetOrder {
		case AscendingOrder:
			d.since = val.Timestamp
		case DescendingOrder:
			d.to = val.Timestamp
		}
	} else {
		val := DatastreamObjectValue{}
		_ = json.Unmarshal([]byte(updateValue.Raw), &val)
		switch d.resultSetOrder {
		case AscendingOrder:
			d.since = val.Timestamp
		case DescendingOrder:
			d.to = val.Timestamp
		}
	}
}

// Parses data obtained by performing a request for a Datastream interface snapshot.
// Returns the snapshot as a map of strings (endpoints) to DatastreamIndividualValues or DatastreamObjectValue,
// depending on the requested interface's aggregation.
func (r GetDatastreamSnapshotResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	return parseDatastreamSnapshot(b, r.aggregation)
}

func parseDatastreamSnapshot(jsonValue []byte, aggregation interfaces.AstarteInterfaceAggregation) (any, error) {
	// clean up useless prefix
	data := gjson.GetBytes(jsonValue, "data")
	if aggregation == interfaces.IndividualAggregation {
		retMap := map[string]any{}
		parseIndividualDatastreamSnapshot([]byte(data.Raw), "", retMap)
		return retMap, nil
	}
	// else, we're dealing with object aggregation (golint is now happy)
	retMap := map[string]DatastreamObjectValue{}
	parseObjectDatastreamSnapshot([]byte(data.Raw), retMap)
	return retMap, nil
}

func parseIndividualDatastreamSnapshot(jsonValue []byte, prefix string, acc map[string]any) {
	// Base case: we have a {"value": n, "timestamp": t} structure
	// a "reception_timestamp" field might also exist, this is handled by unmarshal
	if gjson.GetBytes(jsonValue, "value").Exists() && gjson.GetBytes(jsonValue, "timestamp").Exists() {
		val := DatastreamIndividualValue{}
		_ = json.Unmarshal(jsonValue, &val)
		acc[prefix] = val
		// Recursive case: we have a structure like {"path1": {"value": n, "timestamp": t}, "path2": {"piece2": {"value": n, "timestamp": t}}}
	} else if gjson.ParseBytes(jsonValue).IsObject() {
		insideMap := gjson.ParseBytes(jsonValue).Map()
		for k, v := range insideMap {
			parseIndividualDatastreamSnapshot([]byte(v.Raw), prefix+"/"+k, acc)
		}
	}
	// No third option, maybe we should return an error here
}

func parseObjectDatastreamSnapshot(jsonValue []byte, acc map[string]DatastreamObjectValue) {
	jsonData := gjson.ParseBytes(jsonValue)

	// jsonData must be an object
	obj := jsonData.Value().(map[string]interface{})
	flattened, _ := flat.Flatten(obj, &flat.Options{Safe: true, Delimiter: "."})

	keys := []string{}
	for k := range flattened {
		components := strings.Split(k, ".")
		var theKey string
		if len(components) > 1 {
			theKey = strings.Join(components[:len(components)-1], ".")
		} else {
			theKey = k
		}
		keys = append(keys, theKey)
	}
	keys = removeDuplicateStr(keys)

	rawObjectValues := gjson.GetMany(jsonData.Raw, keys...)
	for i, item := range rawObjectValues {
		value := DatastreamObjectValue{}

		k := fmt.Sprintf("/%s", strings.ReplaceAll(keys[i], ".", "/"))

		if item.IsArray() {
			// since it's a snapshot, we have just one value in the array
			i := item.Array()[0]
			_ = json.Unmarshal([]byte(i.Raw), &value)
			acc[k] = value
		} else {
			_ = json.Unmarshal([]byte(item.Raw), &value)
			acc[k] = value
		}
	}
}

func (r GetDatastreamSnapshotResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request for a property value.
// Returns the value as a PropertyValue.
func (r GetPropertiesResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	// clean up useless prefix
	data := gjson.GetBytes(b, "data")
	retMap := map[string]PropertyValue{}
	parseProperties([]byte(data.Raw), "", retMap)
	return retMap, nil
}

func (r GetPropertiesResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

func parseProperties(jsonValue []byte, prefix string, acc map[string]PropertyValue) {
	// Base case: we have a single value (or an array)
	if !gjson.ParseBytes(jsonValue).IsObject() {
		// leave to the user the choice of type eheh
		acc[prefix] = gjson.ParseBytes(jsonValue).Value()
	} else {
		// Recursive case: we have a structure like {"path2": {"path3": {"path4": n}}}
		insideMap := gjson.ParseBytes(jsonValue).Map()
		for k, v := range insideMap {
			parseProperties([]byte(v.Raw), prefix+"/"+k, acc)
		}
	}
	// No third option, maybe we should return an error here
}

// Parses data obtained by performing a request to list groups for a device.
// Returns the list of groups as an array of strings.
func (r ListGroupsResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data").Array()
	groups := []string{}
	for _, v := range data {
		groups = append(groups, v.Str)
	}
	return groups, nil
}

func (r ListGroupsResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}

// Parses data obtained by performing a request create a new group.
// Returns the group's details as a DevicesAndGroup struct.
func (r CreateGroupResponse) Parse() (any, error) {
	defer r.res.Body.Close()
	b, _ := io.ReadAll(r.res.Body)
	data := gjson.GetBytes(b, "data")
	devicesAndGroup := DevicesAndGroup{}
	_ = json.Unmarshal([]byte(data.Raw), &devicesAndGroup)
	return devicesAndGroup, nil
}

func (r CreateGroupResponse) Raw(f func(*http.Response) any) any {
	defer r.res.Body.Close()
	return f(r.res)
}
