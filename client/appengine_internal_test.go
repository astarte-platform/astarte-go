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
	"reflect"
	"testing"
	"time"

	"github.com/iancoleman/orderedmap"
)

func TestDatastreamSnapshotParsing(t *testing.T) {
	data := `{"test":{"value":2,"timestamp":"2020-03-12T19:46:53.000Z"},"test2":{"value":"2020-03-12T19:51:53.000Z","timestamp":"2020-03-12T19:46:53.000Z","reception_timestamp":"2020-03-12T19:46:53.000Z"},"test3":{"value":"somevalue","timestamp":"2020-03-12T19:46:53.000Z","reception_timestamp":"2020-03-12T19:46:53.346Z"}}`
	v := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		t.Error(err)
	}
	vals, err := parseDatastreamInterface(v)
	if err != nil {
		t.Error(err)
	}

	timestamp, err := time.Parse(time.RFC3339Nano, "2020-03-12T19:46:53.000Z")
	if err != nil {
		t.Error(err)
	}

	if val, ok := vals["/test"]; !ok {
		t.Fail()
	} else {
		if val.Timestamp != timestamp {
			t.Fail()
		}
		if !val.ReceptionTimestamp.IsZero() {
			t.Fail()
		}
	}

	if val, ok := vals["/test2"]; !ok {
		t.Fail()
	} else {
		if val.ReceptionTimestamp != timestamp {
			t.Fail()
		}
	}
}

func TestFailedDatastreamSnapshotParsing(t *testing.T) {
	data := `{"test":[{"value":2,"timestamp":"2020-03-12T19:46:53.000Z"}]}`
	v := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		t.Error(err)
	}
	if vals, err := parseDatastreamInterface(v); err == nil {
		t.Error(vals)
	}
}

func TestAggregateDatastreamSnapshotParsing(t *testing.T) {
	data := `{"test":{"nested":[{"timestamp":"2020-03-23T12:31:08.356Z","val1":true,"val2":12}]}}`

	v := orderedmap.OrderedMap{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		t.Error(err)
	}
	vals, err := parseAggregateDatastreamInterface(v)
	if err != nil {
		t.Error(err)
	}

	timestamp, err := time.Parse(time.RFC3339Nano, "2020-03-23T12:31:08.356Z")
	if err != nil {
		t.Error(err)
	}

	if val, ok := vals["/test/nested"]; !ok {
		t.Error(vals)
	} else {
		if val.Timestamp != timestamp {
			t.Fail()
		}
	}
}

func TestParametricDatastreamParsing(t *testing.T) {
	data := `{"test":{"nested":[{"val1":true,"val2":12,"timestamp":"2022-02-24T11:24:59.284Z"}]}}`

	v := orderedmap.OrderedMap{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		t.Error(err)
	}
	vals, err := parseAggregateDatastreamInterface(v)
	if err != nil {
		t.Error(err)
	}

	timestamp, err := time.Parse(time.RFC3339Nano, "2022-02-24T11:24:59.284Z")
	if err != nil {
		t.Error(err)
	}

	if val, ok := vals["/test/nested"]; !ok {
		t.Error(vals)
	} else {
		if val.Timestamp != timestamp {
			t.Fail()
		}
	}
}

func TestPropertiesParsing(t *testing.T) {
	uglyLookingProperties := "{\"data\":{\"nested\":{\"value\":15, \"timestamp\":\"not_a_timestamp\"}}}"
	// Get the parametric datastream and treat it as individual
	var responseBody struct {
		Data map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal([]byte(uglyLookingProperties), &responseBody)
	if err != nil {
		t.Error(err)
	}

	val := parsePropertyInterface(responseBody.Data)

	if ds, ok := val["/nested/value"]; ok {
		if v, ok := ds.(float64); !ok || v != 15 {
			t.Error("Error in parsing /nested/value", ok, ds, v, reflect.TypeOf(ds))
		}
	} else {
		t.Error("Error in parsing /nested/value", val)
	}

	if ds, ok := val["/nested/timestamp"]; ok {
		if v, ok := ds.(string); !ok || v != "not_a_timestamp" {
			t.Error("Error in parsing /nested/timestamp", ok, ds, v, reflect.TypeOf(ds))
		}
	} else {
		t.Error("Error in parsing /nested/timestamp", val)
	}
}
