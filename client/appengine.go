// Copyright © 2019 Ispirata Srl
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
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/iancoleman/orderedmap"
)

const defaultPageSize int = 10000

var invalidTime time.Time = time.Unix(0, 0)

// AppEngineService is the API Client for AppEngine API
type AppEngineService struct {
	client       *Client
	appEngineURL *url.URL
}

// GetProperties returns all the currently set Properties on a given Interface
func (s *AppEngineService) GetProperties(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	interfaceName string) (map[string]interface{}, error) {
	data, err := s.nestedIndividualQuery(interfaceName, realm, deviceIdentifier, deviceIdentifierType, "")
	if err != nil {
		return nil, err
	}

	return parsePropertyInterface(data), nil
}

// GetDatastreamSnapshot returns all the last values on all paths for a Datastream interface
func (s *AppEngineService) GetDatastreamSnapshot(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	interfaceName string) (map[string]DatastreamValue, error) {
	data, err := s.nestedIndividualQuery(interfaceName, realm, deviceIdentifier, deviceIdentifierType, "")
	if err != nil {
		return nil, err
	}

	return parseDatastreamInterface(data)
}

// GetLastDatastreams returns all the last values on a path for a Datastream interface.
// If limit is <= 0, it returns all existing datastreams. Consider using a GetDatastreamsPaginator in that case.
func (s *AppEngineService) GetLastDatastreams(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, limit int) ([]DatastreamValue, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	return s.getDatastreamInternal(realm, deviceIdentifier, resolvedDeviceIdentifierType, interfaceName, interfacePath, invalidTime, invalidTime, limit, DescendingOrder)
}

// GetDatastreamsPaginator returns a Paginator for all the values on a path for a Datastream interface.
func (s *AppEngineService) GetDatastreamsPaginator(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, resultSetOrder ResultSetOrder) (DatastreamPaginator, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	return s.getDatastreamPaginatorInternal(realm, deviceIdentifier, resolvedDeviceIdentifierType, interfaceName, interfacePath, invalidTime, time.Now(), defaultPageSize, resultSetOrder)
}

// GetDatastreamsTimeWindowPaginator returns a Paginator for all the values on a path in a specified time window for a Datastream interface.
func (s *AppEngineService) GetDatastreamsTimeWindowPaginator(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, since, to time.Time, resultSetOrder ResultSetOrder) (DatastreamPaginator, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	return s.getDatastreamPaginatorInternal(realm, deviceIdentifier, resolvedDeviceIdentifierType, interfaceName, interfacePath, since, to, defaultPageSize, resultSetOrder)
}

// GetAggregateParametricDatastreamSnapshot returns the last value for a Parametric Datastream aggregate interface
func (s *AppEngineService) GetAggregateParametricDatastreamSnapshot(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName string) (map[string]DatastreamAggregateValue, error) {
	// It's a snapshot, so limit=1
	decoder, err := s.appengineGenericJSONDataAPIGet(interfaceName, realm, deviceIdentifier, deviceIdentifierType, "limit=1")
	if err != nil {
		return nil, err
	}
	var responseBody struct {
		Data orderedmap.OrderedMap `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	// If there is no data, return an empty value
	if len(responseBody.Data.Keys()) == 0 {
		return nil, nil
	}

	return parseAggregateDatastreamInterface(responseBody.Data)
}

// GetAggregateDatastreamSnapshot returns the last value for a non-parametric, Datastream aggregate interface
func (s *AppEngineService) GetAggregateDatastreamSnapshot(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName string) (DatastreamAggregateValue, error) {
	// It's a snapshot, so limit=1
	datastreams, err := s.aggregateDatastreamQuery(interfaceName, realm, deviceIdentifier, deviceIdentifierType, "limit=1")
	if err != nil {
		return DatastreamAggregateValue{}, err
	}

	// If there is no data, return an empty value
	if len(datastreams) == 0 {
		return DatastreamAggregateValue{}, nil
	}

	return datastreams[0], nil
}

// GetLastAggregateDatastreams returns the last count values for a Datastream aggregate interface
func (s *AppEngineService) GetLastAggregateDatastreams(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, count int) ([]DatastreamAggregateValue, error) {
	return s.aggregateDatastreamQuery(interfaceName+interfacePath, realm, deviceIdentifier, deviceIdentifierType, fmt.Sprintf("limit=%v", count))
}

// GetAggregateDatastreamsTimeWindow returns the last count values for a Datastream aggregate interface
func (s *AppEngineService) GetAggregateDatastreamsTimeWindow(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, since, to time.Time) ([]DatastreamAggregateValue, error) {
	return s.aggregateDatastreamQuery(interfaceName+interfacePath, realm, deviceIdentifier, deviceIdentifierType,
		fmt.Sprintf("since=%s&to=%s", since.UTC().Format(time.RFC3339Nano), to.UTC().Format(time.RFC3339Nano)))
}

//////////
// Private APIs: These abstract the real calls and do custom decoding of the different reply types
//////////

func (s *AppEngineService) nestedIndividualQuery(urlPath, realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, rawQuery string) (map[string]interface{}, error) {
	decoder, err := s.appengineGenericJSONDataAPIGet(urlPath, realm, deviceIdentifier, deviceIdentifierType, rawQuery)
	if err != nil {
		return nil, err
	}
	var responseBody struct {
		Data map[string]interface{} `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Data, nil
}

func (s *AppEngineService) aggregateDatastreamQuery(urlPath, realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, rawQuery string) ([]DatastreamAggregateValue, error) {
	decoder, err := s.appengineGenericJSONDataAPIGet(urlPath, realm, deviceIdentifier, deviceIdentifierType, rawQuery)
	if err != nil {
		return nil, err
	}
	var responseBody struct {
		Data []DatastreamAggregateValue `json:"data"`
	}
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Data, nil
}

func (s *AppEngineService) appengineGenericJSONDataAPIURL(urlPath, realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, rawQuery string) (*url.URL, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, err := url.Parse(s.appEngineURL.String())
	if err != nil {
		return nil, err
	}
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/interfaces/%s", realm,
		devicePath(deviceIdentifier, resolvedDeviceIdentifierType), urlPath))
	if len(rawQuery) > 0 {
		callURL.RawQuery = rawQuery
	}
	return callURL, nil
}

func (s *AppEngineService) appengineGenericJSONDataAPIGet(urlPath, realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, rawQuery string) (*json.Decoder, error) {
	url, err := s.appengineGenericJSONDataAPIURL(urlPath, realm, deviceIdentifier, deviceIdentifierType, rawQuery)
	if err != nil {
		return nil, err
	}
	return s.client.genericJSONDataAPIGET(url.String(), 200)
}

func (s *AppEngineService) getDatastreamInternal(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string,
	since, to time.Time, limit int, resultSetOrder ResultSetOrder) ([]DatastreamValue, error) {
	realLimit := limit
	if limit < 0 || limit > defaultPageSize {
		realLimit = defaultPageSize
	}
	datastreamPaginator, err := s.getDatastreamPaginatorInternal(realm, deviceIdentifier, deviceIdentifierType, interfaceName, interfacePath,
		since, to, realLimit, resultSetOrder)
	if err != nil {
		return nil, err
	}

	var resultSet []DatastreamValue
	for ok := true; ok; ok = datastreamPaginator.HasNextPage() {
		page, err := datastreamPaginator.GetNextPage()
		if err != nil {
			return nil, err
		}

		// Check special cases
		if limit > 0 {
			totalSize := len(resultSet) + len(page)
			if totalSize == limit {
				return append(resultSet, page...), nil
			} else if totalSize > limit {
				missingSamples := limit - len(resultSet)
				return append(resultSet, page[0:missingSamples-1]...), nil
			}
		}

		resultSet = append(resultSet, page...)
	}

	return resultSet, nil
}

func (s *AppEngineService) getDatastreamPaginatorInternal(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string,
	since, to time.Time, pageSize int, resultSetOrder ResultSetOrder) (DatastreamPaginator, error) {
	url, err := s.appengineGenericJSONDataAPIURL(interfaceName+interfacePath, realm, deviceIdentifier, deviceIdentifierType, "")
	if err != nil {
		return DatastreamPaginator{}, err
	}

	datastreamPaginator := DatastreamPaginator{
		baseURL:        url,
		windowStart:    since,
		windowEnd:      to,
		nextWindow:     invalidTime,
		pageSize:       pageSize,
		client:         s.client,
		hasNextPage:    true,
		resultSetOrder: resultSetOrder,
	}
	return datastreamPaginator, nil
}
