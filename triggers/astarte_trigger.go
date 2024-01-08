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
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type AstarteTriggerMatchOperator string

const (
	All    AstarteTriggerMatchOperator = "*"
	Equal  AstarteTriggerMatchOperator = "=="
	Differ AstarteTriggerMatchOperator = "!="

	Bigger       AstarteTriggerMatchOperator = ">"
	BiggerEqual  AstarteTriggerMatchOperator = ">="
	Smaller      AstarteTriggerMatchOperator = "<"
	SmallerEqual AstarteTriggerMatchOperator = "<="
	Contains     AstarteTriggerMatchOperator = "contains"
	NotContains  AstarteTriggerMatchOperator = "not_contains"
)

// IsValid returns an error if AstarteTriggerType does not represent a valid AstarteTriggerMatchOperator
func (t AstarteTriggerMatchOperator) IsValid() error {
	switch t {
	case All, Equal, Differ, Bigger, BiggerEqual, Smaller, SmallerEqual, Contains, NotContains:
		return nil
	}
	return errors.New("invalid Astarte Trigger type")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *AstarteTriggerMatchOperator) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*t = AstarteTriggerMatchOperator(j)
	if err := t.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid AstarteTriggerMatchOperator", j)
	}
	return nil
}

type AstarteTriggerOn string

const (
	DeviceConnected    AstarteTriggerOn = "device_connected"
	DeviceDisconnected AstarteTriggerOn = "device_disconnected"
	DeviceError        AstarteTriggerOn = "device_error"

	IncomingData       AstarteTriggerOn = "incoming_data"
	ValueStored        AstarteTriggerOn = "value_stored"
	ValueChange        AstarteTriggerOn = "value_change"
	ValueChangeApplied AstarteTriggerOn = "value_change_applied"
	PathCreated        AstarteTriggerOn = "path_created"
	PathRemoved        AstarteTriggerOn = "path_removed"
)

// IsValid returns an error if AstarteTriggerType does not represent a valid AstarteTriggerOn
func (t AstarteTriggerOn) IsValid() error {
	switch t {
	case DeviceConnected, DeviceDisconnected, DeviceError:
		return nil
	case IncomingData, ValueStored, ValueChange, ValueChangeApplied, PathCreated, PathRemoved:
		return nil
	}
	return errors.New("invalid Astarte Trigger type")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *AstarteTriggerOn) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*t = AstarteTriggerOn(j)
	if err := t.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid AstarteTriggerOn", j)
	}
	return nil
}

// AstarteTriggerType represents which kind of Astarte trigger the object represents
type AstarteTriggerType string

const (
	// DataType represents a data trigger
	DataType AstarteTriggerType = "data_trigger"
	// DataType represents a device trigger
	DeviceType AstarteTriggerType = "device_trigger"
)

// IsValid returns an error if AstarteTriggerType does not represent a valid Astarte Trigger Type
func (t AstarteTriggerType) IsValid() error {
	switch t {
	case DataType, DeviceType:
		return nil
	}
	return errors.New("invalid Astarte Trigger type")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *AstarteTriggerType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*t = AstarteTriggerType(j)
	if err := t.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Trigger Type", j)
	}
	return nil
}

// AstarteHTTPMethod represents the kind of http method used
type AstarteHTTPMethod string

const (
	PostMethod   AstarteHTTPMethod = "post"
	GetMethod    AstarteHTTPMethod = "get"
	PutMethod    AstarteHTTPMethod = "put"
	PatchMethod  AstarteHTTPMethod = "patch"
	DeleteMethod AstarteHTTPMethod = "delete"
)

// IsValid returns an error if AstarteHTTPMethod does not represent a valid AstarteHTTPMethod
func (o AstarteHTTPMethod) IsValid() error {
	switch o {
	case PostMethod, GetMethod, PutMethod, PatchMethod, DeleteMethod:
		return nil
	}
	return errors.New("invalid AstarteHTTPMethod")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (o *AstarteHTTPMethod) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*o = AstarteHTTPMethod(j)
	if err := o.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid AstarteHTTPMethod", j)
	}
	return nil
}

type AstarteTriggerAction struct {
	HTTPUrl         string            `json:"http_url"`
	HTTPMethod      AstarteHTTPMethod `json:"http_method"`
	HTTPHeaders     map[string]string `json:"http_static_headers"`
	IgnoreSslErrors bool              `default:"false"`
}
type AstarteSimpleTrigger struct {
	Type               AstarteTriggerType          `json:"type"`
	On                 AstarteTriggerOn            `json:"on"`
	DeviceID           string                      `json:"device_id,omitempty"`
	GroupName          string                      `json:"group_name,omitempty"`
	InterfaceName      string                      `json:"interface_name,omitempty"`
	InterfaceMajor     json.Number                 `json:"interface_major,omitempty"`
	MatchPath          string                      `json:"match_path,omitempty"`
	ValueMatchOperator AstarteTriggerMatchOperator `json:"value_match_operator"`
	KnownValue         *json.Number                `json:"known_value,omitempty"`
}

// AstarteTrigger represents an Astarte Trigger
type AstarteTrigger struct {
	Name           string                 `json:"name"`
	Action         AstarteTriggerAction   `json:"action"`
	SimpleTriggers []AstarteSimpleTrigger `json:"simple_triggers"`
}

// requiredAstarteTrigger is an helper struct used for validating required fields when unmarshalling an
// astarte trigger. Its fields are defined as pointers so that it is possible determining if any field is
// present and valid.
type requiredAstarteTrigger struct {
	Name           *string                        `json:"name"`
	Action         *requiredAstarteTriggerAction  `json:"action"`
	SimpleTriggers []requiredAstarteSimpleTrigger `json:"simple_triggers"`
}
type requiredAstarteTriggerAction struct {
	HTTPUrl    *string            `json:"http_url"`
	HTTPMethod *AstarteHTTPMethod `json:"http_method"`
}

type requiredAstarteSimpleTrigger struct {
	Type      *AstarteTriggerType `json:"type"`
	On        *AstarteTriggerOn   `json:"on"`
	DeviceID  *string             `json:"device_id,omitempty"`
	GroupName *string             `json:"group_name,omitempty"`

	InterfaceName      *string                      `json:"interface_name,omitempty"`
	InterfaceMajor     *json.Number                 `json:"interface_major,omitempty"`
	MatchPath          *string                      `json:"match_path,omitempty"`
	ValueMatchOperator *AstarteTriggerMatchOperator `json:"value_match_operator"`
	KnownValue         *json.Number                 `json:"known_value,omitempty"`
}

// ensureRequiredFields ensures that any required fields within an AstarteTrigger is present and valid. It is
// employed in place of the UnmarshalJSON interface to avoid infinite loops when unmarshalling an AstarteTrigger
//
//nolint:all
func (r *requiredAstarteTrigger) ensureRequiredFields(b []byte) error {

	required := requiredAstarteTrigger{}
	if err := json.Unmarshal(b, &required); err != nil {
		return err
	}
	if required.Name == nil || (required.Name != nil && *required.Name == "") {
		return errors.New("Invalid trigger: name must be set")
	}
	if required.Action == nil {
		return errors.New("Invalid trigger: action must be set")
	}
	if required.Action.HTTPUrl == nil || required.Action.HTTPMethod == nil {
		return errors.New("Invalid trigger: action must have at least an url and a method set")
	}
	if required.Action.HTTPMethod.IsValid() != nil {
		return errors.New("Invalid trigger: invalid method for action")
	}

	if len(required.SimpleTriggers) == 0 {
		return errors.New("Invalid trigger: no triggers are present")
	}
	if len(required.SimpleTriggers) > 1 {
		return errors.New("Invalid trigger: usage of more than one trigger is currently unsupported")
	}

	for _, trigger := range required.SimpleTriggers {
		err2 := simpleTriggerCheck(&trigger)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

//nolint:all
func simpleTriggerCheck(trigger *requiredAstarteSimpleTrigger) error {
	if trigger.Type == nil || trigger.On == nil {
		return errors.New("Invalid trigger condition: Type and On must be set")
	}

	if *trigger.Type != "data_trigger" {

		if *trigger.On != "device_connected" && *trigger.On != "device_disconnected" &&
			*trigger.On != "device_error" {
			return fmt.Errorf("Invalid trigger condition: invalid On value '%v'", *trigger.On)
		}

		if trigger.DeviceID == nil && trigger.GroupName == nil {
			return errors.New("Invalid trigger condition: DeviceID or GroupName must be set")
		}
		if trigger.DeviceID != nil && trigger.GroupName != nil {
			return errors.New("Invalid trigger condition: DeviceID or GroupName cannot both be set ")
		}

		if trigger.InterfaceName != nil ||
			trigger.InterfaceMajor != nil ||
			trigger.MatchPath != nil ||
			trigger.ValueMatchOperator != nil ||
			trigger.KnownValue != nil {
			return errors.New("Invalid trigger: cannot set properties for data trigger on a device trigger")
		}

	} else {
		if *trigger.On != "incoming_data" &&
			*trigger.On != "value_stored" &&
			*trigger.On != "value_change" &&
			*trigger.On != "value_change_applied" &&
			*trigger.On != "path_created" &&
			*trigger.On != "path_removed" {
			return fmt.Errorf("Invalid trigger condition: invalid On value '%v'", *trigger.On)
		}
		if trigger.DeviceID != nil || trigger.GroupName != nil {
			return errors.New("Invalid trigger condition: DeviceID or GroupName cannot be set ")
		}
		if trigger.InterfaceName == nil {
			return errors.New("Invalid data trigger: interface not set, use * to catch all")
		}
		if trigger.InterfaceMajor == nil && *trigger.InterfaceName != "*" {
			return errors.New("Invalid data trigger:  InterfaceMajor must be set")
		}
		if trigger.MatchPath == nil {
			return errors.New("Invalid data trigger: MatchPath not set")
		}
		if trigger.ValueMatchOperator == nil {
			return errors.New("Invalid data trigger: ValueMatchOperator not set")
		}
		if trigger.KnownValue == nil && *trigger.ValueMatchOperator != "*" {
			return errors.New("Invalid data trigger: KnownValue not set")
		}

	}
	return nil
}

// triggerProvider is the object that holds a trigger
type triggerProvider interface {
	[]byte | string
}

// ParseTriggerFrom is a convenience function to call ParseTrigger with an input.
// The input hcan be either a string, tat is interpreted as a file path, or a byteslice.
func ParseTriggerFrom[T triggerProvider](provider T) (AstarteTrigger, error) {
	switch p := any(provider).(type) {
	case string:
		b, err := os.ReadFile(p)
		if err != nil {
			return AstarteTrigger{}, err
		}
		return ParseTrigger(b)
	case []byte:
		return ParseTrigger(p)
	default:
		return AstarteTrigger{}, errors.New("Provided value cannot be used as an Astarte Trigger")
	}
}

// ParseTrigger parses a trigger from a JSON string and returns an AstarteTrigger object when successful.
// Please use this method rather than calling json.Unmarshal on a Trigger, as this will set any missing field
// to the correct, expected default value
func ParseTrigger(triggerContent []byte) (AstarteTrigger, error) {
	astarteTrigger := AstarteTrigger{}
	required := requiredAstarteTrigger{}

	if err := required.ensureRequiredFields(triggerContent); err != nil {
		return astarteTrigger, err
	}

	if err := json.Unmarshal(triggerContent, &astarteTrigger); err != nil {
		return astarteTrigger, err
	}

	return EnsureTriggerDefaults(astarteTrigger), nil
}

// EnsureTriggerDefaults makes sure a JSON-parsed Trigger will have all defaults set. Usually, you should never
// call this method - ParseTrigger does the right thing. It might become useful in case you're dealing with a
// json.Decoder to parse Trigger information
func EnsureTriggerDefaults(astarteTrigger AstarteTrigger) AstarteTrigger {

	// Ensure we have all defaults set
	if err := astarteTrigger.Action.HTTPMethod.IsValid(); err != nil {
		astarteTrigger.Action.HTTPMethod = GetMethod
	}

	subsMapping := []AstarteSimpleTrigger{}
	for _, v := range astarteTrigger.SimpleTriggers {

		if err := v.Type.IsValid(); err != nil {
			v.Type = DataType
		}
		if err := v.On.IsValid(); err != nil {
			v.On = DeviceConnected
		}
		if err := v.ValueMatchOperator.IsValid(); err != nil {
			v.ValueMatchOperator = All
		}
		subsMapping = append(subsMapping, v)
	}
	astarteTrigger.SimpleTriggers = subsMapping

	return astarteTrigger
}
