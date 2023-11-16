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

package triggers

import (
	_ "encoding/json"
	"testing"
)

func TestMissingDataFromTriggerAction(t *testing.T) {

	MissingTriggerName := `
	{
		"action": {
		  "http_url": "https://example.com/my_hook",
		  "http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 0,
			"match_path": "/streamTest/value",
			"value_match_operator": ">",
			"known_value": 0.4
		  }
		]
	  }`

	_, err := ParseTriggerFrom([]byte(MissingTriggerName))
	if err == nil {
		t.Error("This trigger should have failed validation! Missing name")
	}

	EmptyTriggerName := `
	{
		"name": "",
		"action": {
		  "http_url": "https://example.com/my_hook",
		  "http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 0,
			"match_path": "/streamTest/value",
			"value_match_operator": ">",
			"known_value": 0.4
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(EmptyTriggerName))
	if err == nil {
		t.Error("This trigger should have failed validation! Empty name")
	}

	MissinigURL := `
	{
		"name": "test",
		"action": {
		  "http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 0,
			"match_path": "/streamTest/value",
			"value_match_operator": ">",
			"known_value": 0.4
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(MissinigURL))
	if err == nil {
		t.Error("This trigger should have failed validation! Missing URL")
	}

	MissinigMethod := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 0,
			"match_path": "/streamTest/value",
			"value_match_operator": ">",
			"known_value": 0.4
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(MissinigMethod))
	if err == nil {
		t.Error("This trigger should have failed validation! Missing http method")
	}

	DeviceTriggerInvalidHTTP := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "posta"
		},
		"simple_triggers": [
		  {
			"type": "invalid_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidHTTP))
	if err == nil {
		t.Error("This trigger should have Failed! invalid http method")
	}
}

func TestInvalidTriggerData(t *testing.T) {
	DeviceTriggerInvalidMatchOperator := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"",
			"known_value":"3"
		  }
		]
	  }`

	_, err := ParseTriggerFrom([]byte(DeviceTriggerInvalidMatchOperator))
	if err == nil {
		t.Error("This trigger should have Failed! invalid json match operator")
	}

	DeviceTriggerInvalidType := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "invalid_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidType))
	if err == nil {
		t.Error("This trigger should have Failed! invalid trigger type")
	}

	DeviceTriggerInvalidOn := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "AAAAAAAAAAAAAA",
			"device_id": "45336"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidOn))
	if err == nil {
		t.Error("This trigger should have failed validation! invalid 'on' condition")
	}

	DeviceTriggerIncorrectOn := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "incoming_data",
			"device_id": "45336"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerIncorrectOn))
	if err == nil {
		t.Error("This trigger should have failed validation! invalid 'on' condition")
	}

	DeviceTriggerMismatchOn := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerMismatchOn))
	if err == nil {
		t.Error("This trigger should have failed validation! mismatched 'on' condition for trigger type")
	}

}
func TestInvalidTriggerInterface(t *testing.T) {

	DeviceTriggerInterfaceNull := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data"
		  }
		]
	  }`

	_, err := ParseTriggerFrom([]byte(DeviceTriggerInterfaceNull))
	if err == nil {
		t.Error("This trigger should have failed validation! no interfaces specified")
	}

	DeviceTriggerMajorNull := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "AAA"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerMajorNull))
	if err == nil {
		t.Error("This trigger should have failed validation! no interface major specified")
	}

	DeviceTriggerMajorNotNull := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "AAA",
			"interface_major": "2",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerMajorNotNull))
	if err != nil {
		t.Error("This trigger should have passed!")
	}

	DeviceTriggerMajorNullable := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerMajorNullable))
	if err != nil {
		t.Error("This trigger should have passed!")
	}

}

func TestInvalidTriggerGenericErrors(t *testing.T) {

	DeviceTriggerGroupAndID := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "device_connected",
			"device_id": "45336",
			"group_name": "sdgfsd"
		  }
		]
	  }`

	_, err := ParseTriggerFrom([]byte(DeviceTriggerGroupAndID))
	if err == nil {
		t.Error("This trigger should have failed validation! cannot use device_id and group_name")
	}

	DeviceTriggerNoAction := `
	{
		"name": "test",
		"action":"",
		"simple_triggers": [
		  {
			"type": "invalid_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerNoAction))
	if err == nil {
		t.Error("This trigger should have Failed! required action")

	}

	DeviceTriggerNoTriggers := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": []
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerNoTriggers))
	if err == nil {
		t.Error("This trigger should have Failed! one trigger definition is required")
	}

	DeviceTriggerTooManyTriggers := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  },
		  {
			"type": "device_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerTooManyTriggers))
	if err == nil {
		t.Error("This trigger should have Failed! Too many triggers defined")
	}

	DeviceTriggerTypeNotSet := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"==",
			"known_value":"3"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerTypeNotSet))
	if err == nil {
		t.Error("This trigger should have Failed! type not set for trigger")
	}

	DeviceTriggerNoDeviceOrGroup := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "device_connected"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerNoDeviceOrGroup))
	if err == nil {
		t.Error("This trigger should have Failed! device_id or group should be set")
	}

	DeviceTriggerInvalidDevice := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "device_connected",
			"device_id": "dsagfsda",
			"interface_name": "*"

		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidDevice))
	if err == nil {
		t.Error("This trigger should have Failed! invalid data for device")
	}

	DeviceTriggerInvalidDataTrigger := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "device_disconnected"

		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidDataTrigger))
	if err == nil {
		t.Error("This trigger should have Failed! invalid data for trigger type data")
	}

	DeviceTriggerInvalidDataTrigger2 := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"device_id":"34523"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidDataTrigger2))
	if err == nil {
		t.Error("This trigger should have Failed! invalid data for trigger type data")
	}

	DeviceTriggerInvalidDataTrigger3 := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "*"

		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidDataTrigger3))
	if err == nil {
		t.Error("This trigger should have Failed! invalid data for trigger type data")
	}

	DeviceTriggerInvalidDataTrigger4 := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*"

		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidDataTrigger4))
	if err == nil {
		t.Error("This trigger should have Failed! invalid data for trigger type data")
	}

	DeviceTriggerInvalidDataTrigger5 := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path":"/*",
			"value_match_operator":"=="

		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerInvalidDataTrigger5))
	if err == nil {
		t.Error("This trigger should have Failed!invalid data for trigger type data")
	}

}

//nolint:all
func TestParsing(t *testing.T) {

	DataTriggerOk := `
	{
		"name": "example_trigger",
		"action": {
		  "http_url": "https://example.com/my_hook",
		  "http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 0,
			"match_path": "/streamTest/value",
			"value_match_operator": ">",
			"known_value": 0.4
		  }
		]
	  }`

	i, err := ParseTriggerFrom([]byte(DataTriggerOk))
	if err != nil {
		t.Error(err)
	}
	if i.Action.HTTPMethod != PostMethod {
		t.Error("Wrong httpmethod detected", i.Action.HTTPMethod)
	}
	if i.SimpleTriggers[0].Type != DataType {
		t.Error("Wrong type detected", i.SimpleTriggers[0].Type)
	}

	DeviceTriggerOK := `
	{
		"name": "test",
		"action": {
			"http_url": "https://example.com/my_hook",
			"http_method": "post"
		},
		"simple_triggers": [
		  {
			"type": "device_trigger",
			"on": "device_connected",
			"device_id": "45336"
		  }
		]
	  }`

	_, err = ParseTriggerFrom([]byte(DeviceTriggerOK))
	if err != nil {
		t.Error("This trigger should have passed ", err.Error())
	}

}
