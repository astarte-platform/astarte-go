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
	"fmt"
	"strings"
)

// AstarteInterfaceType represents which kind of Astarte interface the object represents
type AstarteInterfaceType int

const (
	// PropertiesType represents a properties Interface
	PropertiesType AstarteInterfaceType = iota
	// DatastreamType represents a datastream Interface
	DatastreamType
)

func (s AstarteInterfaceType) String() string {
	return astarteInterfaceTypeToString[s]
}

var astarteInterfaceTypeToString = map[AstarteInterfaceType]string{
	PropertiesType: "properties",
	DatastreamType: "datastream",
}

var astarteInterfaceTypeToID = map[string]AstarteInterfaceType{
	"properties": PropertiesType,
	"datastream": DatastreamType,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteInterfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteInterfaceTypeToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteInterfaceType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// If the string cannot be found, an error is thrown.
	if val, ok := astarteInterfaceTypeToID[j]; ok {
		*s = val
	} else {
		return fmt.Errorf("'%v' is not a valid Astarte Interface Type", j)
	}
	return nil
}

// AstarteInterfaceOwnership represents the owner of an interface.
type AstarteInterfaceOwnership int

const (
	// DeviceOwnership represents a Device-owned interface
	DeviceOwnership AstarteInterfaceOwnership = iota
	// ServerOwnership represents a Server-owned interface
	ServerOwnership
)

func (s AstarteInterfaceOwnership) String() string {
	return astarteInterfaceOwnershipToString[s]
}

var astarteInterfaceOwnershipToString = map[AstarteInterfaceOwnership]string{
	DeviceOwnership: "device",
	ServerOwnership: "server",
}

var astarteInterfaceOwnershipToID = map[string]AstarteInterfaceOwnership{
	"device": DeviceOwnership,
	"server": ServerOwnership,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteInterfaceOwnership) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteInterfaceOwnershipToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteInterfaceOwnership) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// If the string cannot be found, an error is thrown.
	if val, ok := astarteInterfaceOwnershipToID[j]; ok {
		*s = val
	} else {
		return fmt.Errorf("'%v' is not a valid Astarte Interface Ownership", j)
	}
	return nil
}

// AstarteInterfaceAggregation represents the type of Aggregation of an Interface.
type AstarteInterfaceAggregation int

const (
	// IndividualAggregation represents an interface with individual endpoints
	IndividualAggregation AstarteInterfaceAggregation = iota
	// ObjectAggregation represents an interface with aggregated endpoints
	ObjectAggregation
)

func (s AstarteInterfaceAggregation) String() string {
	return astarteInterfaceAggregationToString[s]
}

var astarteInterfaceAggregationToString = map[AstarteInterfaceAggregation]string{
	IndividualAggregation: "individual",
	ObjectAggregation:     "object",
}

var astarteInterfaceAggregationToID = map[string]AstarteInterfaceAggregation{
	"individual": IndividualAggregation,
	"object":     ObjectAggregation,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteInterfaceAggregation) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteInterfaceAggregationToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteInterfaceAggregation) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'IndividualAggregation' in this case.
	*s = astarteInterfaceAggregationToID[j]
	return nil
}

// AstarteMappingReliability represents the reliability of a mapping
type AstarteMappingReliability int

const (
	// UnreliableReliability represents a QoS 0-like reliability on the wire
	UnreliableReliability AstarteMappingReliability = iota
	// GuaranteedReliability represents a QoS 1-like reliability on the wire
	GuaranteedReliability
	// UniqueReliability represents a QoS 2-like reliability on the wire
	UniqueReliability
)

func (s AstarteMappingReliability) String() string {
	return astarteMappingReliabilityToString[s]
}

var astarteMappingReliabilityToString = map[AstarteMappingReliability]string{
	UnreliableReliability: "unreliable",
	GuaranteedReliability: "guaranteed",
	UniqueReliability:     "unique",
}

var astarteMappingReliabilityToID = map[string]AstarteMappingReliability{
	"unreliable": UnreliableReliability,
	"guaranteed": GuaranteedReliability,
	"unique":     UniqueReliability,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteMappingReliability) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteMappingReliabilityToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteMappingReliability) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'UnreliableReliability' in this case.
	*s = astarteMappingReliabilityToID[j]
	return nil
}

// AstarteMappingRetention represents retention for a single mapping
type AstarteMappingRetention int

const (
	// DiscardRetention means the sample will be discarded if it cannot be sent
	DiscardRetention AstarteMappingRetention = iota
	// VolatileRetention means the sample will be stored in RAM until possible if it cannot be sent
	VolatileRetention
	// StoredRetention means the sample will be stored on Disk until expiration if it cannot be sent
	StoredRetention
)

func (s AstarteMappingRetention) String() string {
	return astarteMappingRetentionToString[s]
}

var astarteMappingRetentionToString = map[AstarteMappingRetention]string{
	DiscardRetention:  "discard",
	VolatileRetention: "volatile",
	StoredRetention:   "stored",
}

var astarteMappingRetentionToID = map[string]AstarteMappingRetention{
	"discard":  DiscardRetention,
	"volatile": VolatileRetention,
	"stored":   StoredRetention,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteMappingRetention) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteMappingRetentionToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteMappingRetention) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'DiscardRetention' in this case.
	*s = astarteMappingRetentionToID[j]
	return nil
}

// AstarteMappingDatabaseRetentionPolicy represents database retention policy for a single mapping
type AstarteMappingDatabaseRetentionPolicy int

const (
	// NoTTL means that there is no expiry (TTL)
	NoTTL AstarteMappingDatabaseRetentionPolicy = iota
	// UseTTL means that database retention TTL is used
	UseTTL
)

func (s AstarteMappingDatabaseRetentionPolicy) String() string {
	return astarteMappingDatabaseRetentionPolicyToString[s]
}

var astarteMappingDatabaseRetentionPolicyToString = map[AstarteMappingDatabaseRetentionPolicy]string{
	NoTTL:  "no_ttl",
	UseTTL: "use_ttl",
}

var astarteMappingDatabaseRetentionPolicyToID = map[string]AstarteMappingDatabaseRetentionPolicy{
	"no_ttl":  NoTTL,
	"use_ttl": UseTTL,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteMappingDatabaseRetentionPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteMappingDatabaseRetentionPolicyToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteMappingDatabaseRetentionPolicy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'NoTTL' in this case.
	*s = astarteMappingDatabaseRetentionPolicyToID[j]
	return nil
}

// AstarteMappingType represents the type of a single mapping. Astarte Types are natively inferred from golang
// native types, as long as the conversion does not lose precision, e.g.: an `int32` value will be accepted
// as a "double" type, but a `float64` value won't be accepted as a "integer" type
type AstarteMappingType int

const (
	// Double represents the "double" type in Astarte. It maps to golang `float64` type,
	// but also accepts implicit conversions from any int or float type
	Double AstarteMappingType = iota
	// Integer represents the "integer" type in Astarte. It maps to golang `int` type,
	// but also accepts implicit conversions from any int type < 64bit
	Integer
	// Boolean represents the "boolean" type in Astarte. It maps to golang `bool` type
	Boolean
	// LongInteger represents the "longinteger" type in Astarte. It maps to golang `int64` type,
	// but also accepts implicit conversions from any int type
	LongInteger
	// String represents the "string" type in Astarte. It maps to golang `string` type
	String
	// BinaryBlob represents the "binaryblob" type in Astarte. It maps to golang `[]byte` type
	BinaryBlob
	// DateTime represents the "datetime" type in Astarte. It maps to golang `time.Time` type
	DateTime
	// DoubleArray represents the "doublearray" type in Astarte. It maps to golang `[]float` type,
	// but also accepts implicit conversions from any int or float type array
	DoubleArray
	// IntegerArray represents the "integerarray" type in Astarte. It maps to golang `[]int` type,
	// but also accepts implicit conversions from any int type < 64bit
	IntegerArray
	// BooleanArray represents the "booleanarray" type in Astarte. It maps to golang `[]bool` type
	BooleanArray
	// LongIntegerArray represents the "longintegerarray" type in Astarte. It maps to golang `[]int64` type,
	// but also accepts implicit conversions from any int type
	LongIntegerArray
	// StringArray represents the "stringarray" type in Astarte. It maps to golang `[]string` type
	StringArray
	// BinaryBlobArray represents the "binaryblobarray" type in Astarte. It maps to golang `[]byte` type
	BinaryBlobArray
	// DateTimeArray represents the "datetimearray" type in Astarte. It maps to golang `[]time.Time` type
	DateTimeArray
)

func (s AstarteMappingType) String() string {
	return astarteMappingTypeToString[s]
}

var astarteMappingTypeToString = map[AstarteMappingType]string{
	Double:           "double",
	Integer:          "integer",
	Boolean:          "boolean",
	LongInteger:      "longinteger",
	String:           "string",
	BinaryBlob:       "binaryblob",
	DateTime:         "datetime",
	DoubleArray:      "doublearray",
	IntegerArray:     "integerarray",
	BooleanArray:     "booleanarray",
	LongIntegerArray: "longintegerarray",
	StringArray:      "stringarray",
	BinaryBlobArray:  "binaryblobarray",
	DateTimeArray:    "datetimearray",
}

var astarteMappingTypeToID = map[string]AstarteMappingType{
	"double":           Double,
	"integer":          Integer,
	"boolean":          Boolean,
	"longinteger":      LongInteger,
	"string":           String,
	"binaryblob":       BinaryBlob,
	"datetime":         DateTime,
	"doublearray":      DoubleArray,
	"integerarray":     IntegerArray,
	"booleanarray":     BooleanArray,
	"longintegerarray": LongIntegerArray,
	"stringarray":      StringArray,
	"binaryblobarray":  BinaryBlobArray,
	"datetimearray":    DateTimeArray,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AstarteMappingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(astarteMappingTypeToString[s])
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AstarteMappingType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Return if found
	if val, ok := astarteMappingTypeToID[j]; ok {
		*s = val
		return nil
	}

	// Fail if not found
	return fmt.Errorf("%s is not a valid Astarte type", j)
}

// AstarteInterfaceMapping represents an individual Mapping in an Astarte Interface
type AstarteInterfaceMapping struct {
	Endpoint                string                                `json:"endpoint"`
	Type                    AstarteMappingType                    `json:"type"`
	Reliability             AstarteMappingReliability             `json:"reliability,omitempty"`
	Retention               AstarteMappingRetention               `json:"retention,omitempty"`
	DatabaseRetentionPolicy AstarteMappingDatabaseRetentionPolicy `json:"database_retention_policy,omitempty"`
	DatabaseRetentionTTL    int                                   `json:"database_retention_ttl,omitempty"`
	Expiry                  int                                   `json:"expiry,omitempty"`
	ExplicitTimestamp       bool                                  `json:"explicit_timestamp,omitempty"`
	AllowUnset              bool                                  `json:"allow_unset,omitempty"`
	Description             string                                `json:"description,omitempty"`
	Documentation           string                                `json:"doc,omitempty"`
}

// AstarteInterface represents an Astarte Interface
type AstarteInterface struct {
	Name              string                      `json:"interface_name"`
	MajorVersion      int                         `json:"version_major"`
	MinorVersion      int                         `json:"version_minor"`
	Type              AstarteInterfaceType        `json:"type"`
	Ownership         AstarteInterfaceOwnership   `json:"ownership"`
	Aggregation       AstarteInterfaceAggregation `json:"aggregation,omitempty"`
	ExplicitTimestamp bool                        `json:"explicit_timestamp,omitempty"`
	HasMetadata       bool                        `json:"has_metadata,omitempty"`
	Description       string                      `json:"description,omitempty"`
	Documentation     string                      `json:"doc,omitempty"`
	Mappings          []AstarteInterfaceMapping   `json:"mappings"`
}

// IsParametric returns whether the interface has at least one parametric mapping
func (a *AstarteInterface) IsParametric() bool {
	for _, v := range a.Mappings {
		if strings.Contains(v.Endpoint, "%{") {
			return true
		}
	}
	return false
}
