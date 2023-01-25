package newclient

import (
	"testing"

	"github.com/astarte-platform/astarte-go/interfaces"
	"github.com/tidwall/gjson"
)

func TestGetDatastreamIndividualSnapshot(t *testing.T) {
	c, _ := getTestContext(t)
	getDatastreamIndividualSnapshotCall, err := c.GetDatastreamIndividualSnapshot(testRealmName, testDeviceID, AstarteDeviceID, testInterfaceName)
	if err != nil {
		t.Error(err)
	}
	res, err := getDatastreamIndividualSnapshotCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	rawData, err := res.Parse()
	if err != nil {
		t.Error(err)
	}
	data, ok := rawData.(map[string]any)
	if !ok {
		t.Errorf("Expected snapshot data map, received %v of type %T", rawData, rawData)
	}
	checkParsedIndividualDatastreamSnapshot(t, data)
}

func TestParseDatastreamIndividualSnapshot(t *testing.T) {
	parsed := map[string]any{}
	parseIndividualDatastreamSnapshot([]byte(gjson.GetBytes([]byte(testIndividualDatastreamSnapshot), "data").Raw), "", parsed)
	checkParsedIndividualDatastreamSnapshot(t, parsed)
}

func TestParseDatastreamObjectSnapshot(t *testing.T) {
	value := `
	{
		"data":{
		   "foo":[
			  {
				 "bar":2,
				 "timestamp":"2022-09-26T14:37:00.468Z",
				 "baz":1
			  }
		   ]
		}
	 }
	`
	retMap := map[string]DatastreamObjectValue{}
	parseObjectDatastreamSnapshot([]byte(gjson.GetBytes([]byte(value), "data").Raw), "", retMap)
	for k, v := range retMap {
		if k == "/foo" {
			barV, ok := v.Values.Get("bar")
			if !ok {
				t.Errorf("Value not found: bar")
			}
			bar := barV.(float64)
			bazV, ok := v.Values.Get("baz")
			if !ok {
				t.Errorf("Value not found: baz")
			}
			baz := bazV.(float64)
			if !(bar == 2 && baz == 1) {
				t.Errorf("Unexpected values: bar %v , baz: %v\n", bar, baz)
			}
		} else {
			t.Error("Unexpected path")
		}
	}
}

func TestParseDatastreamObject(t *testing.T) {
	value := `
	{
		"data":
		[
			{
				"bar":1,
				"timestamp":"2022-09-26T13:38:22.627Z",
				"baz":0
			},
			{
				"bar":2,
				"timestamp":"2022-09-26T14:37:00.468Z",
				"baz":1
			}
		]
	}
	`
	parsed := []DatastreamObjectValue{}
	parseDatastream([]byte(gjson.GetBytes([]byte(value), "data").Raw), "")
	for _, v := range parsed {
		barV, ok := v.Values.Get("bar")
		if !ok {
			t.Errorf("Value not found: bar")
		}
		bar := barV.(float64)
		bazV, ok := v.Values.Get("baz")
		if !ok {
			t.Errorf("Value not found: baz")
		}
		baz := bazV.(float64)
		if !(bar == 2 && baz == 1) && !(bar == 1 && baz == 0) {
			t.Errorf("Unexpected values: bar %v , baz: %v\n", bar, baz)
		}
	}
}
func TestParseProperties(t *testing.T) {
	value := `
	{
		"data":{
		   "their":{
			  "new":{
				 "value":11
			  }
		   }
		}
	 }
	`
	retMap := map[string]PropertyValue{}
	parseProperties([]byte(gjson.GetBytes([]byte(value), "data").Raw), "", retMap)
	for k, v := range retMap {
		if k == "/their/new/value" {
			value := v.(float64)
			if value != 11 {
				t.Errorf("Unexpected value: %v of type %T\n", v, v)
			}
		} else {
			t.Error("Unexpected path")
		}
	}
}

func TestSendData(t *testing.T) {
	simpleMapping := interfaces.AstarteInterfaceMapping{Endpoint: "/an/endpoint", Type: interfaces.Integer, AllowUnset: true}
	datastreamInterface := interfaces.AstarteInterface{Name: testServerOwnedInterfaceName, Ownership: interfaces.ServerOwnership, Type: interfaces.DatastreamType, Mappings: []interfaces.AstarteInterfaceMapping{simpleMapping}, Aggregation: interfaces.IndividualAggregation}
	parametricMapping := interfaces.AstarteInterfaceMapping{Endpoint: "/%{an}/endpoint", Type: interfaces.Integer, AllowUnset: true}
	parametricInterface := interfaces.AstarteInterface{Name: testServerOwnedInterfaceName, Ownership: interfaces.ServerOwnership, Type: interfaces.DatastreamType, Mappings: []interfaces.AstarteInterfaceMapping{parametricMapping}, Aggregation: interfaces.IndividualAggregation}
	propertyInterface := interfaces.AstarteInterface{Name: testServerOwnedPropertyInterfaceName, Ownership: interfaces.ServerOwnership, Type: interfaces.PropertiesType, Mappings: []interfaces.AstarteInterfaceMapping{simpleMapping}}

	c, _ := getTestContext(t)
	sendDatastreamCall, err := c.SendData(testRealmName, testDeviceID, AstarteDeviceID, datastreamInterface, "/an/endpoint", 42)
	if err != nil {
		t.Error(err)
	}
	res, err := sendDatastreamCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}

	sendParametricDatastreamCall, err := c.SendData(testRealmName, testDeviceID, AstarteDeviceID, parametricInterface, "/other/endpoint", 42)
	if err != nil {
		t.Error(err)
	}
	res, err = sendParametricDatastreamCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}

	setPropertyCall, err := c.SendData(testRealmName, testDeviceID, AstarteDeviceID, propertyInterface, "/an/endpoint", 42)
	if err != nil {
		t.Error(err)
	}
	res, err = setPropertyCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}

	unsetPropertyCall, err := c.UnsetProperty(testRealmName, testDeviceID, AstarteDeviceID, testServerOwnedPropertyInterfaceName, "/an/endpoint")
	if err != nil {
		t.Error(err)
	}
	res, err = unsetPropertyCall.Run(c)
	if err != nil {
		t.Error(err)
	}
	_, err = res.Parse()
	if err != nil {
		t.Error(err)
	}
}

func checkParsedIndividualDatastreamSnapshot(t *testing.T, result map[string]any) {
	for k, v := range result {
		if k == "/anotherTest/value" {
			value := v.(DatastreamIndividualValue)
			if value.Value != 0.29031942518908505 {
				t.Error("Unexpected value")
			}
		} else if k == "/yetAnotherTest/value" {
			value := v.(DatastreamIndividualValue)
			if value.Value != 0.41505074846327805 {
				t.Error("Unexpected value")
			}
		} else {
			t.Error("Unexpected path")
		}
	}
}
