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

// DeviceResultFormat represents the format of the Device returned in the Device list.
type DeviceResultFormat int

const (
	// DeviceIDFormat means the Paginator will return a list of strings
	// representing the Device ID of the Devices.
	DeviceIDFormat DeviceResultFormat = iota
	// DeviceDetailsFormat means the Paginator will return a list of
	// DeviceDetails structs
	DeviceDetailsFormat
)

// DeviceListPaginator handles a paginated set of results. It provides a one-directional iterator to call onto
// Astarte AppEngine API and handle potentially extremely large sets of results in chunk. You should prefer
// DeviceListPaginator rather than direct API calls if you expect your result set to be particularly large.
type DeviceListPaginator struct {
	baseURL     *url.URL
	nextQuery   url.Values
	format      DeviceResultFormat
	pageSize    int
	client      *Client
	hasNextPage bool
}

// Rewind rewinds the simulator to the first page. GetNextPage will then return the first page of the call.
func (d *DeviceListPaginator) Rewind() {
	d.nextQuery = url.Values{}
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

// GetNextPage retrieves the next result page from the paginator and populates
// the array pointed by pagePtr with it.
// The type of pagePtr must be the correct one depending on the format of the
// paginator: if format is DeviceIDFormat pagePtr must be *[]string, if format
// is DeviceDetailsFormat pagePtr must be *[]DeviceDetails.
// If no more results are available, HasNextPage will return false. GetNextPage
// throws an error if no more pages are available.
func (d *DeviceListPaginator) GetNextPage(pagePtr interface{}) error {
	if !d.hasNextPage {
		return errors.New("No more pages available")
	}

	if err := d.checkPageFormat(pagePtr); err != nil {
		return err
	}

	callURL, _ := d.setupCallURL()

	links := Links{}
	err := d.client.genericJSONDataAPIGETWithLinks(pagePtr, &links, callURL.String(), 200)
	if err != nil {
		return err
	}

	d.computePageState(&links)

	return nil
}

func (d *DeviceListPaginator) checkPageFormat(pagePtr interface{}) error {
	switch d.format {
	case DeviceIDFormat:
		_, ok := pagePtr.(*[]string)
		if !ok {
			return errors.New("pagePtr must be of type *[]string when using DeviceIDFormat")
		}

	case DeviceDetailsFormat:
		_, ok := pagePtr.(*[]DeviceDetails)
		if !ok {
			return errors.New("pagePtr must be of type *[]DeviceDetails when using DeviceDetailsFormat")
		}
	}

	return nil
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

	query := d.nextQuery
	switch d.format {
	case DeviceIDFormat:
		query.Set("details", "false")
	case DeviceDetailsFormat:
		query.Set("details", "true")
	}

	callURL.RawQuery = query.Encode()

	return callURL, nil
}
