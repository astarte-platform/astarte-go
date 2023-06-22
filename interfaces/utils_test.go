// Copyright © 2020-2023 SECO Mind Srl
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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
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

	// Validate queries
	if err := ValidateQuery(i, "/testSensor/name"); err != nil {
		t.Error(err)
	}
	if err := ValidateQuery(i, "/testSensor"); err != nil {
		t.Error(err)
	}
	if err := ValidateQuery(i, "/"); err != nil {
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

	// Validate queries
	if err := ValidateQuery(i, "/testSensor/names"); err == nil {
		t.Fail()
	}
	if err := ValidateQuery(i, "/testSensor/names/extra"); err == nil {
		t.Fail()
	}
}

func TestAggregateMessageValidation(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "datastream",
		"ownership": "device",
		"aggregation": "object",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/sensors/%{sensor_id}/name",
				"type": "string",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name.",
				"retention": "discard",
				"reliability": "unreliable",
				"database_retention_policy": "use_ttl",
				"database_retention_ttl": 200
			},
			{
				"endpoint": "/sensors/%{sensor_id}/unit",
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

	if err := ValidateInterfacePath(i, "/sensors/testSensor/name"); err != nil {
		t.Error(err)
	}
	if err := ValidateAggregateMessage(i, "/sensors/testSensor", map[string]interface{}{"name": "test"}); err != nil {
		t.Error(err)
	}

	// Validate queries
	if err := ValidateQuery(i, "/sensors/testSensor"); err != nil {
		t.Error(err)
	}
	if err := ValidateQuery(i, "/sensors"); err != nil {
		t.Error(err)
	}
	if err := ValidateQuery(i, "/"); err != nil {
		t.Error(err)
	}
}

func TestAggregateMessageWrongPaths(t *testing.T) {
	validInterface := `
	{
		"interface_name": "org.astarte-platform.genericsensors.AvailableSensors",
		"version_major": 0,
		"version_minor": 1,
		"type": "datastream",
		"ownership": "device",
		"aggregation": "object",
		"description": "Describes available generic sensors.",
		"doc": "This interface allows to describe available sensors and their attributes such as name and sampled data measurement unit. Sensors are identified by their sensor_id. See also org.astarte-platform.genericsensors.AvailableSensors.",
		"mappings": [
			{
				"endpoint": "/sensors/%{sensor_id}/name",
				"type": "string",
				"description": "Sensor name.",
				"doc": "An arbitrary sensor name.",
				"retention": "discard",
				"reliability": "unreliable",
				"database_retention_policy": "use_ttl",
				"database_retention_ttl": 200
			},
			{
				"endpoint": "/sensors/%{sensor_id}/unit",
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

	if err := ValidateInterfacePath(i, "/sensors/testSensor/name/extra"); err == nil {
		t.Fail()
	}
	if err := ValidateInterfacePath(i, "/sensors/testSensor/names"); err == nil {
		t.Fail()
	}
	if err := ValidateAggregateMessage(i, "/sensors/testSensor", map[string]interface{}{"names": "check"}); err == nil {
		t.Fail()
	}
	if _, err := InterfaceMappingFromPath(i, "/sensors/testSensor/names"); err == nil {
		t.Fail()
	}

	// Validate queries
	if err := ValidateQuery(i, "/sensors/testSensor/name"); err == nil {
		t.Fail()
	}
	if err := ValidateQuery(i, "/sensorsa/testSensor"); err == nil {
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
	if err := ValidateIndividualMessage(i, "/binaryblobarrayValue", [][]byte{{'b', 'l', 'o', 'b'}}); err != nil {
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

func TestPayloadNormalization(t *testing.T) {
	byteArray := []byte{'a', 's', 't', 'a', 'r', 't', 'e'}
	if NormalizePayload(byteArray, true).(string) != base64.StdEncoding.EncodeToString(byteArray) {
		t.Error("Base64 matching in normalization failed")
	}
	if !bytes.Equal(NormalizePayload(byteArray, false).([]byte), byteArray) {
		t.Error("Base64 matching in normalization failed")
	}

	timestamp := time.Now()
	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		t.Error(err)
	}

	if NormalizePayload(timestamp.In(loc), true) != timestamp.UTC() {
		t.Error("Time conversion failed", timestamp.In(loc), timestamp.UTC())
	}

	inMap := map[string]time.Time{"testTime": timestamp.In(loc)}
	outMap := map[string]interface{}{"testTime": timestamp.UTC()}

	if !reflect.DeepEqual(NormalizePayload(inMap, true), outMap) {
		t.Error("Map conversion failed", NormalizePayload(inMap, true), outMap)
	}

	inInterfaceMap := map[string]interface{}{"testTime": timestamp.In(loc)}

	if !reflect.DeepEqual(NormalizePayload(inInterfaceMap, true), outMap) {
		t.Error("Map conversion failed", NormalizePayload(inInterfaceMap, true), outMap)
	}

	inSlice := [][]byte{byteArray}
	outSlice := []string{base64.StdEncoding.EncodeToString(byteArray)}

	if !reflect.DeepEqual(NormalizePayload(inSlice, true), outSlice) {
		t.Error("Slice conversion failed", NormalizePayload(inSlice, true), outSlice)
	}

	inInterfaceSlice := []interface{}{byteArray}
	outInterfaceSlice := []interface{}{base64.StdEncoding.EncodeToString(byteArray)}

	if !reflect.DeepEqual(NormalizePayload(inInterfaceSlice, true), outInterfaceSlice) {
		t.Error("Slice conversion failed", NormalizePayload(inInterfaceSlice, true), outInterfaceSlice)
	}

	inMultiMap := map[string]interface{}{
		"testTime":            timestamp.In(loc),
		"testBytearray":       byteArray,
		"testNestedBytearray": [][]byte{byteArray},
		"testString":          "test",
		"testStringArray":     []string{"test"},
	}
	outMultiMapEncoded := map[string]interface{}{
		"testTime":            timestamp.UTC(),
		"testBytearray":       base64.StdEncoding.EncodeToString(byteArray),
		"testNestedBytearray": []string{base64.StdEncoding.EncodeToString(byteArray)},
		"testString":          "test",
		"testStringArray":     []interface{}{"test"},
	}
	outMultiMapNonEncoded := map[string]interface{}{
		"testTime":            timestamp.UTC(),
		"testBytearray":       byteArray,
		"testNestedBytearray": [][]byte{byteArray},
		"testString":          "test",
		"testStringArray":     []interface{}{"test"},
	}

	if !reflect.DeepEqual(NormalizePayload(inMultiMap, true), outMultiMapEncoded) {
		t.Error("Multimap conversion failed", NormalizePayload(inMultiMap, true), outMultiMapEncoded)
	}
	if !reflect.DeepEqual(NormalizePayload(inMultiMap, false), outMultiMapNonEncoded) {
		t.Error("Multimap conversion failed", NormalizePayload(inMultiMap, false), outMultiMapNonEncoded)
	}
}
