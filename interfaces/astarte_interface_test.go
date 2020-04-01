// Copyright © 2020 Ispirata Srl
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

package interfaces

import (
	"encoding/json"
	"testing"
)

func TestParsing(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/name",
				"type": "string",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name.",
				"retention": "discard",
				"reliability": "unreliable",
				"database_retention_policy": "use_ttl",
				"database_retention_ttl": 200
			},
			{
				"endpoint": "/%{sensor_id}/unit",
				"type": "string",
				"description": "Sample data measurement unit.",
				"doc": "SI unit such as m, kg, K, etc..."
			}
		]
	}`

	i, err := ParseInterfaceFromString(validInterface)
	if err != nil {
		t.Error(err)
	}
	if i.Aggregation != IndividualAggregation {
		t.Error("Wrong aggregation detected", i.Aggregation)
	}
	if i.Mappings[0].Retention != DiscardRetention {
		t.Error("Wrong retention detected", i.Mappings[0].Retention)
	}
	if i.Mappings[1].Retention != DiscardRetention {
		t.Error("Wrong retention detected", i.Mappings[0].Retention)
	}
}

func TestMarshaling(t *testing.T) {
	i := AstarteInterface{
		Name:          "org.astarte-platform.genericsensors.AvailableSensors",
		MajorVersion:  1,
		MinorVersion:  0,
		Type:          PropertiesType,
		Ownership:     DeviceOwnership,
		Description:   "Describes available generic sensors.",
		Documentation: "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		Mappings: []AstarteInterfaceMapping{
			AstarteInterfaceMapping{
				Endpoint:                "/%{sensor_id}/name",
				Type:                    String,
				Description:             "Sensor name.",
				Documentation:           "An arbitrary sensor name.",
				Retention:               DiscardRetention,
				Reliability:             UnreliableReliability,
				DatabaseRetentionPolicy: UseTTL,
				DatabaseRetentionTTL:    30000,
			},
			AstarteInterfaceMapping{
				Endpoint:      "/%{sensor_id}/unit",
				Type:          String,
				Description:   "Sample data measurement unit.",
				Documentation: "SI unit such as m, kg, K, etc...",
			},
		},
	}

	if _, err := json.Marshal(i); err != nil {
		t.Error(err)
	}
}

func TestFailedTypeParsing(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/name",
				"type": "stringa",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name."
			},
			{
				"endpoint": "/%{sensor_id}/unit",
				"type": "string",
				"description": "Sample data measurement unit.",
				"doc": "SI unit such as m, kg, K, etc..."
			}
		]
	}`

	if _, err := ParseInterfaceFromString(validInterface); err == nil {
		t.Error("This interface should have failed validation!")
	}
}

func TestFailedStructureParsing(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/name",
				"type": "strings",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name."
			},
			{
				"endpoint": "/%{sensor_id}/unit",
				"description": "Sample data measurement unit.",
				"doc": "SI unit such as m, kg, K, etc..."
			}
		]
	}`

	if _, err := ParseInterfaceFromString(validInterface); err == nil {
		t.Error("This interface should have failed validation!")
	}
}

func TestFailedMarshalingParsing(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_minor": "test",
		"type": 3,
		"ownership": 2,
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/name",
				"type": "strings",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name."
			},
			{
				"endpoint": "/%{sensor_id}/unit",
				"description": "Sample data measurement unit.",
				"doc": "SI unit such as m, kg, K, etc..."
			}
		]
	}`

	if _, err := ParseInterfaceFromString(validInterface); err == nil {
		t.Error("This interface should have failed validation!")
	}
}

func TestFailedOwnershipParsing(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "devices",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/%{sensor_id}/name",
				"type": "string",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name."
			},
			{
				"endpoint": "/%{sensor_id}/unit",
				"type": "string",
				"description": "Sample data measurement unit.",
				"doc": "SI unit such as m, kg, K, etc..."
			}
		]
	}`

	if _, err := ParseInterfaceFromString(validInterface); err == nil {
		t.Error(err)
	}
}
