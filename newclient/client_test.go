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
	"testing"
)

func TestClientValidation(t *testing.T) {
	// A standard client
	if _, err := New(WithBaseURL("api.an-astarte.org")); err != nil {
		t.Error(err)
	}

	// Client with conflicting URLs
	if _, err := New(
		WithBaseURL("api.an-astarte.org"),
		WithAppengineURL("a.different.appengine-url.com"),
		WithUserAgent("pippo"),
	); err == nil {
		t.Error("Conflicting URLs were given to client, but no error found")
	}

	// Invalid URL provided
	if _, err := New(
		WithBaseURL("://api.an-astarte.org/thethings"),
	); err == nil {
		t.Error("Invalid URL provided, but no error found")
	}
}
