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
	"time"
)

func TestMessageValidation(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"aggregation": "individual",
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
				"database_retention_policy": "ttl",
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

	i := AstarteInterface{}
	if err := json.Unmarshal([]byte(validInterface), &i); err != nil {
		t.Error(err)
	}

	if err := ValidateInterfacePath(i, "/testSensor/name"); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/testSensor/name", "test"); err != nil {
		t.Error(err)
	}
}

func TestParametricMessageWrongPaths(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"aggregation": "individual",
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
				"database_retention_policy": "ttl",
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

	i := AstarteInterface{}
	if err := json.Unmarshal([]byte(validInterface), &i); err != nil {
		t.Error(err)
	}

	if err := ValidateInterfacePath(i, "/testSensor/name/extra"); err == nil {
		t.Fail()
	}
	if err := ValidateInterfacePath(i, "/testSensor/names"); err == nil {
		t.Fail()
	}
	if err := ValidateIndividualMessage(i, "/testSensor/names", "check"); err == nil {
		t.Fail()
	}
	if _, err := InterfaceMappingFromPath(i, "/testSensor/extra/path"); err == nil {
		t.Fail()
	}
}

func TestAggregateMessageValidation(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"aggregation": "object",
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
				"database_retention_policy": "ttl",
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

	i := AstarteInterface{}
	if err := json.Unmarshal([]byte(validInterface), &i); err != nil {
		t.Error(err)
	}

	if err := ValidateInterfacePath(i, "/testSensor/name"); err != nil {
		t.Error(err)
	}
	if err := ValidateAggregateMessage(i, map[string]interface{}{"/testSensor/name": "test"}); err != nil {
		t.Error(err)
	}
}

func TestAggregateMessageWrongPaths(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "properties",
		"ownership": "device",
		"aggregation": "object",
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
				"database_retention_policy": "ttl",
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

	i := AstarteInterface{}
	if err := json.Unmarshal([]byte(validInterface), &i); err != nil {
		t.Error(err)
	}

	if err := ValidateInterfacePath(i, "/testSensor/name/extra"); err == nil {
		t.Fail()
	}
	if err := ValidateInterfacePath(i, "/testSensor/names"); err == nil {
		t.Fail()
	}
	if err := ValidateAggregateMessage(i, map[string]interface{}{"/testSensor/names": "check"}); err == nil {
		t.Fail()
	}
	if _, err := InterfaceMappingFromPath(i, "/testSensor/names"); err == nil {
		t.Fail()
	}
}

func TestTypeValidation(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.tests.TypeValidation",
		"version_major": 0,
		"version_minor": 1,
		"type": "datastream",
		"ownership": "device",
		"mappings": [
			{
				"endpoint": "/integerValue",
				"type": "integer"
			},
			{
				"endpoint": "/longintegerValue",
				"type": "longinteger"
			},
			{
				"endpoint": "/doubleValue",
				"type": "double"
			},
			{
				"endpoint": "/stringValue",
				"type": "string"
			},
			{
				"endpoint": "/booleanValue",
				"type": "boolean"
			},
			{
				"endpoint": "/binaryblobValue",
				"type": "binaryblob"
			},
			{
				"endpoint": "/datetimeValue",
				"type": "datetime"
			},
			{
				"endpoint": "/integerarrayValue",
				"type": "integerarray"
			},
			{
				"endpoint": "/longintegerarrayValue",
				"type": "longintegerarray"
			},
			{
				"endpoint": "/doublearrayValue",
				"type": "doublearray"
			},
			{
				"endpoint": "/stringarrayValue",
				"type": "stringarray"
			},
			{
				"endpoint": "/booleanarrayValue",
				"type": "booleanarray"
			},
			{
				"endpoint": "/binaryblobarrayValue",
				"type": "binaryblobarray"
			},
			{
				"endpoint": "/datetimearrayValue",
				"type": "datetimearray"
			}
		]
	}`

	i := AstarteInterface{}
	if err := json.Unmarshal([]byte(validInterface), &i); err != nil {
		t.Error(err)
	}

	if err := ValidateIndividualMessage(i, "/integerValue", 42); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/integerValue", int32(42)); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/integerValue", int16(42)); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/integerValue", int8(42)); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/longintegerValue", 42); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/longintegerValue", int64(42)); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/doubleValue", 3.14); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/doubleValue", 314); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/stringValue", "test"); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/booleanValue", true); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/binaryblobValue", []byte{'b', 'l', 'o', 'b'}); err != nil {
		t.Error(err)
	}
	timestamp := time.Now()
	if err := ValidateIndividualMessage(i, "/datetimeValue", timestamp); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/datetimeValue", &timestamp); err != nil {
		t.Error(err)
	}

	// Arrays
	if err := ValidateIndividualMessage(i, "/integerarrayValue", []int{42}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/integerarrayValue", []int32{int32(42)}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/integerarrayValue", []int16{int16(42)}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/integerarrayValue", []int8{int8(42)}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/longintegerarrayValue", []int{42}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/longintegerarrayValue", []int64{int64(42)}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/doublearrayValue", []float32{3.14}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/doublearrayValue", []float64{3.14}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/doublearrayValue", []int{314}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/stringarrayValue", []string{"test"}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/booleanarrayValue", []bool{true}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/binaryblobarrayValue", [][]byte{[]byte{'b', 'l', 'o', 'b'}}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/datetimearrayValue", []time.Time{timestamp}); err != nil {
		t.Error(err)
	}
	if err := ValidateIndividualMessage(i, "/datetimearrayValue", []*time.Time{&timestamp}); err != nil {
		t.Error(err)
	}
}

func TestFailedTypeValidation(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.tests.TypeValidation",
		"version_major": 0,
		"version_minor": 1,
		"type": "datastream",
		"ownership": "device",
		"mappings": [
			{
				"endpoint": "/integerValue",
				"type": "integer"
			},
			{
				"endpoint": "/longintegerValue",
				"type": "longinteger"
			},
			{
				"endpoint": "/doubleValue",
				"type": "double"
			},
			{
				"endpoint": "/stringValue",
				"type": "string"
			},
			{
				"endpoint": "/booleanValue",
				"type": "boolean"
			},
			{
				"endpoint": "/binaryblobValue",
				"type": "binaryblob"
			},
			{
				"endpoint": "/datetimeValue",
				"type": "datetime"
			},
			{
				"endpoint": "/integerarrayValue",
				"type": "integerarray"
			},
			{
				"endpoint": "/longintegerarrayValue",
				"type": "longintegerarray"
			},
			{
				"endpoint": "/doublearrayValue",
				"type": "doublearray"
			},
			{
				"endpoint": "/stringarrayValue",
				"type": "stringarray"
			},
			{
				"endpoint": "/booleanarrayValue",
				"type": "booleanarray"
			},
			{
				"endpoint": "/binaryblobarrayValue",
				"type": "binaryblobarray"
			},
			{
				"endpoint": "/datetimearrayValue",
				"type": "datetimearray"
			}
		]
	}`

	i := AstarteInterface{}
	if err := json.Unmarshal([]byte(validInterface), &i); err != nil {
		t.Error(err)
	}

	if err := ValidateIndividualMessage(i, "/integerValue", "test"); err == nil {
		t.Fail()
	}
	if err := ValidateIndividualMessage(i, "/integerValue", int64(42)); err == nil {
		t.Fail()
	}
	if err := ValidateIndividualMessage(i, "/integerValue", 4.2); err == nil {
		t.Fail()
	}
	if err := ValidateIndividualMessage(i, "/longintegerValue", 4.2); err == nil {
		t.Fail()
	}
	if err := ValidateIndividualMessage(i, "/stringValue", 42); err == nil {
		t.Fail()
	}
}
