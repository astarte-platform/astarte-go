// Copyright © 2023 SECO Mind Srl
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
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrConflictingUrls               = errors.New("Conflicting URLs provided")
	ErrNoUrlsProvided                = errors.New("No Astarte URL(s) provided")
	ErrNoPrivateKeyProvided          = errors.New("No Astarte private key provided")
	ErrRealmNameNotProvided          = errors.New("Realm name was not provided")
	ErrRealmPublicKeyNotProvided     = errors.New("Realm public key was not provided")
	ErrTooManyReplicationFactors     = errors.New("Can't have both replication factor and datacenter replication factors")
	ErrNegativeReplicationFactor     = errors.New("Replication factor must be a strictly positive integer")
	ErrTooHighExpiry                 = errors.New("Expiry for tokens generated from a private key must be less than 5 minutes")
	ErrNoAuthProvided                = errors.New("Neither an Astarte JWT nor an Astarte private key were provided")
	ErrBothJWTAndPrivateKey          = errors.New("Can't provide both an Astarte JWT and an Astarte private key")
	ErrExpiryButNoPrivateKeyProvided = errors.New("Expiry was set, but no Astarte private key provided")
)

func ErrInvalidDeviceID(deviceID string) error {
	return fmt.Errorf("%s is not a valid Astarte device ID", deviceID)
}

func ErrDifferentStatusCode(expected, received int) error {
	return fmt.Errorf("Received unexpeced status code: %d instead of %d", received, expected)
}

func errorFromJSONErrors(responseBody io.Reader) error {
	var errorBody struct {
		Errors map[string]interface{} `json:"errors"`
	}

	err := json.NewDecoder(responseBody).Decode(&errorBody)
	if err != nil {
		return err
	}

	errJSON, _ := json.MarshalIndent(&errorBody, "", "  ")
	return fmt.Errorf("%s", errJSON)
}

func runAstarteRequestError(res *http.Response, expectedCode int) (AstarteResponse, error) {
	if res.Body != nil {
		return Empty{}, errorFromJSONErrors(res.Body)
	}
	return Empty{}, ErrDifferentStatusCode(expectedCode, res.StatusCode)
}
