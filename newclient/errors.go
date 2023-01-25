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

package newclient

import (
	"errors"
	"fmt"
)

var (
	ErrConflictingUrls           error = errors.New("Conflicting URLs provided")
	ErrNoUrlsProvided            error = errors.New("No Astarte URL(s) provided")
	ErrNoPrivateKeyProvided      error = errors.New("No Astarte private key provided")
	ErrDifferentStatusCode       error = errors.New("Received unexpected status code")
	ErrRealmNameNotProvided      error = errors.New("Realm name was not provided")
	ErrRealmPublicKeyNotProvided error = errors.New("Realm public key was not provided")
	ErrTooManyReplicationFactors error = errors.New("Can't have both replication factor and datacenter replication factors")
	ErrNegativeReplicationFactor error = errors.New("Replication factor must be a strictly positive integer")
)

func ErrInvalidDeviceID(deviceID string) error {
	return errors.New(fmt.Sprintf("%s is not a valid Astarte device ID.", deviceID))
}
