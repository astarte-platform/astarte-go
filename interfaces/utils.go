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
)

// ValidateAggregateMessage validates an aggregate message
func ValidateAggregateMessage(astarteInterface AstarteInterface, values map[string]interface{}) error {
	for k, v := range values {
		// Validate the individual message
		if err := ValidateIndividualMessage(astarteInterface, "/"+k, v); err != nil {
			return err
		}
		// TODO: validate the type
	}

	return nil
}

// ValidateIndividualMessage validates an individual message
func ValidateIndividualMessage(astarteInterface AstarteInterface, path string, value interface{}) error {
	// Get the corresponding mapping
	_, err := InterfaceMappingFromPath(astarteInterface, path)
	if err != nil {
		return err
	}

	// TODO: validate the type

	return nil
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
