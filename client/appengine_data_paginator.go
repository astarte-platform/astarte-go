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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/astarte-platform/astarte-go/interfaces"
	"moul.io/http2curl"
)

// ResultSetOrder represents the order of the samples.
type ResultSetOrder int

const (
	// AscendingOrder means the Paginator will return results starting from the oldest.
	AscendingOrder ResultSetOrder = iota
	// DescendingOrder means the Paginator will return results starting from the newest.
	DescendingOrder
)

// DatastreamPaginator handles a paginated set of results. It provides a one-directional iterator to call onto
// Astarte AppEngine API and handle potentially extremely large sets of results in chunk.
type DatastreamPaginator struct {
	baseURL        *url.URL
	since          time.Time
	to             time.Time
	firstPage      bool
	nextQuery      url.Values
	resultSetOrder ResultSetOrder
	pageSize       int
	client         *Client
	hasNextPage    bool
	aggregation    interfaces.AstarteInterfaceAggregation
}

// Rewind rewinds the paginator to the first page. GetNextPage will then return the first page of the call.
func (d *DatastreamPaginator) Rewind() {
	// Invalid time
	d.since = time.Time{}
	d.to = time.Time{}
	d.nextQuery = url.Values{}
	d.hasNextPage = true
	d.firstPage = true
}

// HasNextPage returns whether this paginator can return more pages.
func (d *DatastreamPaginator) HasNextPage() bool {
	return d.hasNextPage
}

// GetPageSize returns the page size for this paginator.
func (d *DatastreamPaginator) GetPageSize() int {
	return d.pageSize
}

// GetResultSetOrder returns the order in which samples are returned for this paginator.
func (d *DatastreamPaginator) GetResultSetOrder() ResultSetOrder {
	return d.resultSetOrder
}

// GetNextPage returns a request to get the next result page from the paginator.
// If no more results are available, HasNextPage will return false.
// GetNextPage throws an error if no more pages are available or if an invalid parameter is specified.
func (d *DatastreamPaginator) GetNextPage() (AstarteRequest, error) {
	if !d.hasNextPage {
		return nil, errors.New("No more pages available")
	}

	callURL, err := d.setupCallURL()
	if err != nil {
		return Empty{}, err
	}
	req := d.client.makeHTTPrequest(http.MethodGet, callURL, nil)

	return GetNextDatastreamPageRequest{req: req, expects: 200, paginator: d}, nil
}

type GetNextDatastreamPageRequest struct {
	req       *http.Request
	expects   int
	paginator Paginator
}

// nolint:bodyclose
func (r GetNextDatastreamPageRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return r.handleNextDatastreamPageFail(res)
	}
	return GetNextDatastreamPageResponse{res: res, paginator: &r.paginator}, nil
}

func (r GetNextDatastreamPageRequest) handleNextDatastreamPageFail(res *http.Response) (AstarteResponse, error) {
	if res.Body == nil {
		return Empty{}, ErrDifferentStatusCode(r.expects, res.StatusCode)
	}
	// A quirky corner case:
	// when the size of Astarte data is a multiple of r.paginator.pageSize,
	// the last page will be too far in the future and the last request will fail.
	// Let's make sure everything works correctly.
	p, _ := r.paginator.(*DatastreamPaginator)
	if !p.firstPage {
		return GetNextDatastreamPageResponse{res: res, paginator: &r.paginator}, nil
	}
	// now that the corner case is handled, if we're here we must fail
	return Empty{}, errorFromJSONErrors(res.Body)
}

func (r GetNextDatastreamPageRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

func (d *DatastreamPaginator) setupCallURL() (*url.URL, error) {
	callURL, _ := url.Parse(d.baseURL.String())

	query := d.nextQuery
	switch d.resultSetOrder {
	case AscendingOrder:
		// If no start is set, let's start from the beginnning of time
		if (d.since == time.Time{}) {
			d.since = time.Unix(0, 0)
		}
		// All data in the next page come from a time after 'since' (so we descend)
		if d.firstPage {
			// first page includes also the starting value
			query.Set("since", d.since.UTC().Format(time.RFC3339Nano))
		} else {
			// pages after the first must not include the starting value
			query.Set("since_after", d.since.UTC().Format(time.RFC3339Nano))
			query.Del("since")
		}
		if (d.to != time.Time{}) {
			// All data in the next page come from a time until 'to'
			query.Set("to", d.to.UTC().Format(time.RFC3339Nano))
		}
		if d.pageSize != 0 {
			query.Set("limit", fmt.Sprintf("%d", d.pageSize))
		}

	case DescendingOrder:
		if d.pageSize == 0 {
			return &url.URL{}, fmt.Errorf("A limit parameter must be specified when using DescendingOrder")
		}
		if (d.since != time.Time{}) {
			return &url.URL{}, fmt.Errorf("A since parameter must not be specified when using DescendingOrder")
		}
		query.Set("limit", fmt.Sprintf("%d", d.pageSize))
		// if "to" doesn't exist, default behavior with only "limit" is descending
		if (d.to != time.Time{}) {
			// All data in the next page come from a time until 'to' (so we descend)
			query.Set("to", d.to.UTC().Format(time.RFC3339Nano))
		}
	}

	callURL.RawQuery = query.Encode()

	return callURL, nil
}
