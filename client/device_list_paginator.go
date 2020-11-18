// Copyright © 2019-2020 Ispirata Srl
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
	"errors"
	"net/url"
)

// DeviceListPaginator handles a paginated set of results. It provides a one-directional iterator to call onto
// Astarte AppEngine API and handle potentially extremely large sets of results in chunk. You should prefer
// DeviceListPaginator rather than direct API calls if you expect your result set to be particularly large.
type DeviceListPaginator struct {
	baseURL     *url.URL
	nextQuery   url.Values
	pageSize    int
	client      *Client
	hasNextPage bool
}

// Rewind rewinds the simulator to the first page. GetNextPage will then return the first page of the call.
func (d *DeviceListPaginator) Rewind() {
	// We remove `from_token` query parameter
	d.nextQuery.Del("from_token")
	d.hasNextPage = true
}

// HasNextPage returns whether this paginator can return more pages
func (d *DeviceListPaginator) HasNextPage() bool {
	return d.hasNextPage
}

// GetPageSize returns the page size for this paginator
func (d *DeviceListPaginator) GetPageSize() int {
	return d.pageSize
}

// GetNextPage retrieves the next result page from the paginator. Returns the
// page as an array of Device IDs. If no more results are available,
// HasNextPage will return false. GetNextPage throws an error if no more pages
// are available.
func (d *DeviceListPaginator) GetNextPage() ([]string, error) {
	if !d.hasNextPage {
		return nil, errors.New("No more pages available")
	}

	callURL, _ := d.setupCallURL()

	page := []string{}
	links := Links{}
	err := d.client.genericJSONDataAPIGETWithLinks(&page, &links, callURL.String(), 200)
	if err != nil {
		return nil, err
	}

	d.computePageState(&links)

	return page, nil
}

func (d *DeviceListPaginator) computePageState(links *Links) {
	if links.Next == "" {
		d.hasNextPage = false
	} else {
		d.hasNextPage = true
		parsedLinks, _ := url.Parse(links.Next)
		d.nextQuery = parsedLinks.Query()
	}
}

func (d *DeviceListPaginator) setupCallURL() (*url.URL, error) {
	callURL, err := url.Parse(d.baseURL.String())
	if err != nil {
		return nil, err
	}
	callURL.RawQuery = d.nextQuery.Encode()

	return callURL, nil
}
