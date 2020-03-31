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

	// Note that if the string is empty then it will be set to 'IndividualAggregation'
	if j == "" {
		*a = IndividualAggregation
	} else {
		*a = AstarteInterfaceAggregation(j)
		if err := a.IsValid(); err != nil {
			return fmt.Errorf("'%v' is not a valid Astarte Interface Aggregation", j)
		}
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

	// Note that if the string is empty then it will be set to 'UnreliableReliability'
	if j == "" {
		*r = UnreliableReliability
	} else {
		*r = AstarteMappingReliability(j)
		if err := r.IsValid(); err != nil {
			return fmt.Errorf("'%v' is not a valid Astarte Mapping Reliability", j)
		}
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

	// Note that if the string is empty then it will be set to 'DiscardRetention'
	if j == "" {
		*r = DiscardRetention
	} else {
		*r = AstarteMappingRetention(j)
		if err := r.IsValid(); err != nil {
			return fmt.Errorf("'%v' is not a valid Astarte Mapping Retention", j)
		}
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

	// Note that if the string is empty then it will be set to 'NoTTL'
	if j == "" {
		*r = NoTTL
	} else {
		*r = AstarteMappingDatabaseRetentionPolicy(j)
		if err := r.IsValid(); err != nil {
			return fmt.Errorf("'%v' is not a valid Astarte Mapping Database Retention Policy", j)
		}
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

// IsParametric returns whether the interface has at least one parametric mapping
func (a *AstarteInterface) IsParametric() bool {
	for _, v := range a.Mappings {
		if strings.Contains(v.Endpoint, "%{") {
			return true
		}
	}
	return false
}
