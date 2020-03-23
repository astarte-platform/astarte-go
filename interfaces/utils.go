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
	"fmt"
	"strings"
	"time"
)

// ValidateAggregateMessage validates an aggregate message
func ValidateAggregateMessage(astarteInterface AstarteInterface, values map[string]interface{}) error {
	for k, v := range values {
		// Validate the individual message
		if err := ValidateIndividualMessage(astarteInterface, "/"+k, v); err != nil {
			return err
		}
	}

	return nil
}

// ValidateIndividualMessage validates an individual message
func ValidateIndividualMessage(astarteInterface AstarteInterface, path string, value interface{}) error {
	// Get the corresponding mapping
	mapping, err := InterfaceMappingFromPath(astarteInterface, path)
	if err != nil {
		return err
	}

	// Validate type and return result
	return validateType(mapping.Type, value)
}

// InterfaceMappingFromPath retrieves the corresponding interface mapping given a path, and returns a meaningful error
// the path cannot be resolved.
func InterfaceMappingFromPath(astarteInterface AstarteInterface, interfacePath string) (AstarteInterfaceMapping, error) {
	// Ensure we're matching exactly one of the mappings.
	if !astarteInterface.IsParametric() {
		return simpleMappingValidation(astarteInterface, interfacePath)
	}

	return parametricMappingValidation(astarteInterface, interfacePath)

}

// ValidateInterfacePath validates path against the structure of astarteInterface, and returns a meaningful error
// the path cannot be resolved.
func ValidateInterfacePath(astarteInterface AstarteInterface, interfacePath string) error {
	_, err := InterfaceMappingFromPath(astarteInterface, interfacePath)
	return err
}

func simpleMappingValidation(astarteInterface AstarteInterface, interfacePath string) (AstarteInterfaceMapping, error) {
	// Is the path valid?
	for _, mapping := range astarteInterface.Mappings {
		if mapping.Endpoint == interfacePath {
			return mapping, nil
		}
	}
	return AstarteInterfaceMapping{}, fmt.Errorf("Path %s does not exist on Interface %s", interfacePath, astarteInterface.Name)
}

func parametricMappingValidation(astarteInterface AstarteInterface, interfacePath string) (AstarteInterfaceMapping, error) {
	// Is the path valid?
	interfacePathTokens := strings.Split(interfacePath, "/")
	for _, mapping := range astarteInterface.Mappings {
		mappingTokens := strings.Split(mapping.Endpoint, "/")
		if len(mappingTokens) != len(interfacePathTokens) {
			continue
		}
		// Iterate
		matchFound := true
		for index, token := range mappingTokens {
			if interfacePathTokens[index] != token && !strings.HasPrefix(token, "%{") {
				matchFound = false
				break
			}
		}
		if matchFound {
			return mapping, nil
		}
	}
	return AstarteInterfaceMapping{}, fmt.Errorf("Path %s does not exist on Interface %s", interfacePath, astarteInterface.Name)
}

func validateType(mappingType AstarteMappingType, value interface{}) error {
	// Do a case switch and check, depending on the golang type of value, whether
	// we have a match with the Astarte type or not.
	switch value.(type) {
	case int, int8, int16, int32, uint, uint16, uint32:
		if mappingType == Integer || mappingType == LongInteger || mappingType == Double {
			return nil
		}
	case int64, uint64:
		if mappingType == LongInteger || mappingType == Double {
			return nil
		}
	case float32, float64:
		if mappingType == Double {
			return nil
		}
	case string:
		if mappingType == String {
			return nil
		}
	case bool:
		if mappingType == Boolean {
			return nil
		}
	case []byte:
		if mappingType == BinaryBlob {
			return nil
		}
	case time.Time, *time.Time:
		if mappingType == DateTime {
			return nil
		}
	case []int, []int8, []int16, []int32, []uint, []uint16, []uint32:
		if mappingType == IntegerArray || mappingType == LongIntegerArray || mappingType == DoubleArray {
			return nil
		}
	case []int64:
		if mappingType == LongIntegerArray || mappingType == DoubleArray {
			return nil
		}
	case []float32, []float64:
		if mappingType == DoubleArray {
			return nil
		}
	case []string:
		if mappingType == StringArray {
			return nil
		}
	case []bool:
		if mappingType == BooleanArray {
			return nil
		}
	case [][]byte:
		if mappingType == BinaryBlobArray {
			return nil
		}
	case []time.Time, []*time.Time:
		if mappingType == DateTimeArray {
			return nil
		}
	}

	return fmt.Errorf("Value for mapping does not match type restrictions for %s", mappingType.String())
}
