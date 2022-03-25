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
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// AstarteInterfaceType represents which kind of Astarte interface the object represents
type AstarteInterfaceType string

const (
	// PropertiesType represents a properties Interface
	PropertiesType AstarteInterfaceType = "properties"
	// DatastreamType represents a datastream Interface
	DatastreamType AstarteInterfaceType = "datastream"
)

// IsValid returns an error if AstarteInterfaceType does not represent a valid Astarte Interface Type
func (t AstarteInterfaceType) IsValid() error {
	switch t {
	case PropertiesType, DatastreamType:
		return nil
	}
	return errors.New("invalid Astarte Interface type")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *AstarteInterfaceType) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*t = AstarteInterfaceType(j)
	if err := t.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Interface Type", j)
	}
	return nil
}

// AstarteInterfaceOwnership represents the owner of an interface.
type AstarteInterfaceOwnership string

const (
	// DeviceOwnership represents a Device-owned interface
	DeviceOwnership AstarteInterfaceOwnership = "device"
	// ServerOwnership represents a Server-owned interface
	ServerOwnership AstarteInterfaceOwnership = "server"
)

// IsValid returns an error if AstarteInterfaceOwnership does not represent a valid Astarte Ownership Type
func (o AstarteInterfaceOwnership) IsValid() error {
	switch o {
	case DeviceOwnership, ServerOwnership:
		return nil
	}
	return errors.New("invalid Astarte Interface ownership")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (o *AstarteInterfaceOwnership) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*o = AstarteInterfaceOwnership(j)
	if err := o.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Interface Ownership", j)
	}
	return nil
}

// AstarteInterfaceAggregation represents the type of Aggregation of an Interface.
type AstarteInterfaceAggregation string

const (
	// IndividualAggregation represents an interface with individual endpoints
	IndividualAggregation AstarteInterfaceAggregation = "individual"
	// ObjectAggregation represents an interface with aggregated endpoints
	ObjectAggregation AstarteInterfaceAggregation = "object"
)

// IsValid returns an error if AstarteInterfaceAggregation does not represent a valid Astarte Interface Aggregation
func (a AstarteInterfaceAggregation) IsValid() error {
	switch a {
	case IndividualAggregation, ObjectAggregation:
		return nil
	}
	return errors.New("invalid Astarte Interface aggregation")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (a *AstarteInterfaceAggregation) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*a = AstarteInterfaceAggregation(j)
	if err := a.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Interface Aggregation", j)
	}
	return nil
}

// AstarteMappingReliability represents the reliability of a mapping
type AstarteMappingReliability string

const (
	// UnreliableReliability represents a QoS 0-like reliability on the wire
	UnreliableReliability AstarteMappingReliability = "unreliable"
	// GuaranteedReliability represents a QoS 1-like reliability on the wire
	GuaranteedReliability AstarteMappingReliability = "guaranteed"
	// UniqueReliability represents a QoS 2-like reliability on the wire
	UniqueReliability AstarteMappingReliability = "unique"
)

// IsValid returns an error if AstarteMappingReliability does not represent a valid Astarte Mapping Reliability
func (r AstarteMappingReliability) IsValid() error {
	switch r {
	case UnreliableReliability, GuaranteedReliability, UniqueReliability:
		return nil
	}
	return errors.New("invalid Astarte Mapping reliability")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *AstarteMappingReliability) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*r = AstarteMappingReliability(j)
	if err := r.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Mapping Reliability", j)
	}
	return nil
}

// AstarteMappingRetention represents retention for a single mapping
type AstarteMappingRetention string

const (
	// DiscardRetention means the sample will be discarded if it cannot be sent
	DiscardRetention AstarteMappingRetention = "discard"
	// VolatileRetention means the sample will be stored in RAM until possible if it cannot be sent
	VolatileRetention AstarteMappingRetention = "volatile"
	// StoredRetention means the sample will be stored on Disk until expiration if it cannot be sent
	StoredRetention AstarteMappingRetention = "stored"
)

// IsValid returns an error if AstarteMappingRetention does not represent a valid Astarte Mapping Retention
func (r AstarteMappingRetention) IsValid() error {
	switch r {
	case DiscardRetention, VolatileRetention, StoredRetention:
		return nil
	}
	return errors.New("invalid Astarte Mapping retention")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *AstarteMappingRetention) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*r = AstarteMappingRetention(j)
	if err := r.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Mapping Retention", j)
	}
	return nil
}

// AstarteMappingDatabaseRetentionPolicy represents database retention policy for a single mapping
type AstarteMappingDatabaseRetentionPolicy string

const (
	// NoTTL means that there is no expiry (TTL)
	NoTTL AstarteMappingDatabaseRetentionPolicy = "no_ttl"
	// UseTTL means that database retention TTL is used
	UseTTL AstarteMappingDatabaseRetentionPolicy = "use_ttl"
)

// IsValid returns an error if AstarteMappingDatabaseRetentionPolicy does not represent a valid Astarte Mapping Database Retention Policy
func (r AstarteMappingDatabaseRetentionPolicy) IsValid() error {
	switch r {
	case NoTTL, UseTTL:
		return nil
	}
	return errors.New("invalid Astarte Mapping database retention policy")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (r *AstarteMappingDatabaseRetentionPolicy) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*r = AstarteMappingDatabaseRetentionPolicy(j)
	if err := r.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Mapping Database Retention Policy", j)
	}
	return nil
}

// AstarteMappingType represents the type of a single mapping. Astarte Types are natively inferred from golang
// native types, as long as the conversion does not lose precision, e.g.: an `int32` value will be accepted
// as a "double" type, but a `float64` value won't be accepted as a "integer" type
type AstarteMappingType string

const (
	// Double represents the "double" type in Astarte. It maps to golang `float64` type,
	// but also accepts implicit conversions from any int or float type
	Double AstarteMappingType = "double"
	// Integer represents the "integer" type in Astarte. It maps to golang `int` type,
	// but also accepts implicit conversions from any int type < 64bit
	Integer AstarteMappingType = "integer"
	// Boolean represents the "boolean" type in Astarte. It maps to golang `bool` type
	Boolean AstarteMappingType = "boolean"
	// LongInteger represents the "longinteger" type in Astarte. It maps to golang `int64` type,
	// but also accepts implicit conversions from any int type
	LongInteger AstarteMappingType = "longinteger"
	// String represents the "string" type in Astarte. It maps to golang `string` type
	String AstarteMappingType = "string"
	// BinaryBlob represents the "binaryblob" type in Astarte. It maps to golang `[]byte` type
	BinaryBlob AstarteMappingType = "binaryblob"
	// DateTime represents the "datetime" type in Astarte. It maps to golang `time.Time` type
	DateTime AstarteMappingType = "datetime"
	// DoubleArray represents the "doublearray" type in Astarte. It maps to golang `[]float` type,
	// but also accepts implicit conversions from any int or float type array
	DoubleArray AstarteMappingType = "doublearray"
	// IntegerArray represents the "integerarray" type in Astarte. It maps to golang `[]int` type,
	// but also accepts implicit conversions from any int type < 64bit
	IntegerArray AstarteMappingType = "integerarray"
	// BooleanArray represents the "booleanarray" type in Astarte. It maps to golang `[]bool` type
	BooleanArray AstarteMappingType = "booleanarray"
	// LongIntegerArray represents the "longintegerarray" type in Astarte. It maps to golang `[]int64` type,
	// but also accepts implicit conversions from any int type
	LongIntegerArray AstarteMappingType = "longintegerarray"
	// StringArray represents the "stringarray" type in Astarte. It maps to golang `[]string` type
	StringArray AstarteMappingType = "stringarray"
	// BinaryBlobArray represents the "binaryblobarray" type in Astarte. It maps to golang `[]byte` type
	BinaryBlobArray AstarteMappingType = "binaryblobarray"
	// DateTimeArray represents the "datetimearray" type in Astarte. It maps to golang `[]time.Time` type
	DateTimeArray AstarteMappingType = "datetimearray"
)

// IsValid returns an error if AstarteMappingType does not represent a valid Astarte Mapping Type
func (m AstarteMappingType) IsValid() error {
	switch m {
	case Double, Integer, Boolean, LongInteger, String, BinaryBlob, DateTime,
		DoubleArray, IntegerArray, BooleanArray, LongIntegerArray, StringArray, BinaryBlobArray, DateTimeArray:
		return nil
	}
	return errors.New("invalid Astarte Mapping type")
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (m *AstarteMappingType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*m = AstarteMappingType(j)
	if err := m.IsValid(); err != nil {
		return fmt.Errorf("'%v' is not a valid Astarte Mapping Type", j)
	}
	return nil
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

// requiredAstarteInterface is an helper struct used for validating required fields when unmarshalling an
// astarte interface. Its fields are defined as pointers so that it is possible determining if any field is
// present and valid.
type requiredAstarteInterface struct {
	Name         *string                           `json:"interface_name"`
	MajorVersion *int                              `json:"version_major"`
	MinorVersion *int                              `json:"version_minor"`
	Type         *string                           `json:"type"`
	Ownership    *string                           `json:"ownership"`
	Mappings     []requiredAstarteInterfaceMapping `json:"mappings"`
}

// requiredAstarteInterfaceMapping is an helper struct used for validating required fields when unmarshalling an
// astarte interface mapping. Its fields are defined as pointers so that it is possible determining if any field is
// present and valid.
type requiredAstarteInterfaceMapping struct {
	Endpoint *string `json:"endpoint"`
	Type     *string `json:"type"`
}

// ensureRequiredFields ensures that any required fields within an AstarteInterface is present and valid. It is
// employed in place of the UnmarshalJSON interface to avoid infinite loops when unmarshalling an AstarteInterface
func (r *requiredAstarteInterface) ensureRequiredFields(b []byte) error {
	required := requiredAstarteInterface{}
	if err := json.Unmarshal(b, &required); err != nil {
		return err
	}
	if required.Name == nil || (required.Name != nil && *required.Name == "") {
		return errors.New("Invalid interface: interface_name must be set")
	}
	if required.MajorVersion == nil {
		return errors.New("Invalid interface: version_major must be set")
	}
	if required.MinorVersion == nil {
		return errors.New("Invalid interface: version_minor must be set")
	}
	if required.Type == nil {
		return errors.New("Invalid interface: type must be set")
	}
	if required.Ownership == nil {
		return errors.New("Invalid interface: ownership must be set")
	}
	if len(required.Mappings) == 0 {
		return errors.New("Invalid interface: no mappings are present")
	}
	for _, m := range required.Mappings {
		if m.Endpoint == nil || (m.Endpoint != nil && *m.Endpoint == "") {
			return errors.New("Invalid interface: missing endpoint in mapping")
		}
		if m.Type == nil {
			return errors.New("Invalid interface: missing type in mapping")
		}
	}
	return nil
}

// ParseInterfaceFromFile is a convenience function to call ParseInterface with a file as input
func ParseInterfaceFromFile(interfaceFile string) (AstarteInterface, error) {
	b, err := ioutil.ReadFile(interfaceFile)
	if err != nil {
		return AstarteInterface{}, err
	}
	return ParseInterface(b)
}

// ParseInterfaceFromString is a convenience function to call ParseInterface with a string as input
func ParseInterfaceFromString(interfaceContent string) (AstarteInterface, error) {
	return ParseInterface([]byte(interfaceContent))
}

// ParseInterface parses an interface from a JSON string and returns an AstarteInterface object when successful.
// Please use this method rather than calling json.Unmarshal on an interface, as this will set any missing field
// to the correct, expected default value
func ParseInterface(interfaceContent []byte) (AstarteInterface, error) {
	astarteInterface := AstarteInterface{}
	required := requiredAstarteInterface{}

	if err := required.ensureRequiredFields(interfaceContent); err != nil {
		return astarteInterface, err
	}

	if err := json.Unmarshal(interfaceContent, &astarteInterface); err != nil {
		return astarteInterface, err
	}

	return EnsureInterfaceDefaults(astarteInterface), nil
}

// EnsureInterfaceDefaults makes sure a JSON-parsed interface will have all defaults set. Usually, you should never
// call this method - ParseInterface does the right thing. It might become useful in case you're dealing with a
// json.Decoder to parse interface information
func EnsureInterfaceDefaults(astarteInterface AstarteInterface) AstarteInterface {
	// Ensure we have all defaults set
	if err := astarteInterface.Aggregation.IsValid(); err != nil {
		astarteInterface.Aggregation = IndividualAggregation
	}

	subsMapping := []AstarteInterfaceMapping{}
	for _, v := range astarteInterface.Mappings {
		if err := v.Reliability.IsValid(); err != nil {
			v.Reliability = UnreliableReliability
		}
		if err := v.Retention.IsValid(); err != nil {
			v.Retention = DiscardRetention
		}
		if err := v.DatabaseRetentionPolicy.IsValid(); err != nil {
			v.DatabaseRetentionPolicy = NoTTL
		}
		subsMapping = append(subsMapping, v)
	}
	astarteInterface.Mappings = subsMapping

	return astarteInterface
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
