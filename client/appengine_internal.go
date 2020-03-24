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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astarte-platform/astarte-go/misc"
	"github.com/iancoleman/orderedmap"
)

// resolveDeviceIdentifierType maps a deviceIdentifier and DeviceIdentifierType to a resolved
// DeviceIdentifierType (i.e. AstarteDeviceID or AstarteDeviceAlias). AutodiscoverDeviceIdentifier
// is resolved by checking if it's a valid Device ID, otherwise it's considered a Device Alias.
// AstarteDeviceAlias and AstarteDeviceID are returned as is.
func resolveDeviceIdentifierType(deviceIdentifier string, deviceIdentifierType DeviceIdentifierType) DeviceIdentifierType {
	switch deviceIdentifierType {
	case AutodiscoverDeviceIdentifier:
		if misc.IsValidAstarteDeviceID(deviceIdentifier) {
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

////////
// Parsing functions - mostly useful to parse snapshots e.g.: API call to the interface itself.
////////

func parsePropertyInterface(interfaceMap map[string]interface{}) map[string]interface{} {
	// Start recursion and return resulting map
	return parsePropertiesMap(interfaceMap, "")
}

func parseDatastreamInterface(interfaceMap map[string]interface{}) (map[string]DatastreamValue, error) {
	// Start recursion and return resulting map
	return parseDatastreamMap(interfaceMap, "")
}

func parseAggregateDatastreamInterface(interfaceMap orderedmap.OrderedMap) (map[string]DatastreamAggregateValue, error) {
	// Start recursion and return resulting map
	return parseAggregateDatastreamMap(interfaceMap, "")
}

func parsePropertiesMap(aMap map[string]interface{}, completeKeyPath string) map[string]interface{} {
	m := make(map[string]interface{})

	for key, val := range aMap {
		switch actualVal := val.(type) {
		case map[string]interface{}:
			for k, v := range parsePropertiesMap(actualVal, completeKeyPath+"/"+key) {
				m[k] = v
			}
		default:
			m[completeKeyPath+"/"+key] = actualVal
		}
	}

	return m
}

func parseAggregateDatastreamMap(aMap orderedmap.OrderedMap, completeKeyPath string) (map[string]DatastreamAggregateValue, error) {
	m := make(map[string]DatastreamAggregateValue)

	// Special case: have we hit the bottom?
	if val, ok := aMap.Get("timestamp"); ok {
		// Corner case: this might actually be just a token in the path named "timestamp". Let's ensure it
		// does not contain an object.
		if _, ok := val.(map[string]interface{}); !ok {
			datastreamValue, err := parseAggregateDatastreamValue(aMap)
			if err != nil {
				return nil, err
			}
			m[completeKeyPath] = datastreamValue
			return m, nil
		}
	}

	foundAnything := false
	for _, key := range aMap.Keys() {
		val, _ := aMap.Get(key)
		switch cVal := val.(type) {
		case orderedmap.OrderedMap:
			foundAnything = true
			parsedMap, err := parseAggregateDatastreamMap(cVal, completeKeyPath+"/"+key)
			if err != nil {
				return nil, err
			}
			for k, v := range parsedMap {
				m[k] = v
			}
		}
	}
	if !foundAnything {
		return m, errors.New("Could not parse Datastream - payload is likely malformed")
	}

	return m, nil
}

func parseAggregateDatastreamValue(aMap orderedmap.OrderedMap) (DatastreamAggregateValue, error) {
	// Ensure some type safety
	var timestamp time.Time
	timestampInterface, _ := aMap.Get("timestamp")
	switch t := timestampInterface.(type) {
	case time.Time:
		timestamp = t
	case string:
		var err error
		timestamp, err = time.Parse(time.RFC3339Nano, t)
		if err != nil {
			return DatastreamAggregateValue{}, err
		}
	}

	aMap.Delete("timestamp")
	return DatastreamAggregateValue{Values: aMap, Timestamp: timestamp}, nil
}

func parseDatastreamMap(aMap map[string]interface{}, completeKeyPath string) (map[string]DatastreamValue, error) {
	m := make(map[string]DatastreamValue)

	// Special case: have we hit the bottom?
	if val, ok := aMap["value"]; ok {
		// Corner case: this might actually be just a token in the path named "value". Let's ensure it
		// does not contain an object.
		if _, ok := val.(map[string]interface{}); !ok {
			datastreamValue, err := parseDatastreamValue(aMap)
			if err != nil {
				return nil, err
			}
			m[completeKeyPath] = datastreamValue
			return m, nil
		}
	}

	foundAnything := false
	for key, val := range aMap {
		switch cVal := val.(type) {
		case map[string]interface{}:
			foundAnything = true
			parsedMap, err := parseDatastreamMap(cVal, completeKeyPath+"/"+key)
			if err != nil {
				return nil, err
			}
			for k, v := range parsedMap {
				m[k] = v
			}
		}
	}
	if !foundAnything {
		return m, errors.New("Could not parse Datastream - payload is likely malformed")
	}

	return m, nil
}

func parseDatastreamValue(aMap map[string]interface{}) (DatastreamValue, error) {
	// Ensure some type safety
	// Unmarshal into DatastreamValue
	d := DatastreamValue{}
	jsonData, err := json.Marshal(aMap)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return d, err
	}

	return d, nil
}
