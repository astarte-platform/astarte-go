package newclient

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"moul.io/http2curl"
)

// ResultSetOrder represents the order of the samples.
type ResultSetOrder int

const (
	// AscendingOrder means the Paginator will return results starting from the oldest.
	AscendingOrder ResultSetOrder = iota
	// DescendingOrder means the Paginator will return results starting from the oldest.
	DescendingOrder
)

// DatastreamPaginator handles a paginated set of results. It provides a one-directional iterator to call onto
// Astarte AppEngine API and handle potentially extremely large sets of results in chunk.
type DatastreamPaginator struct {
	baseURL              *url.URL
	windowOlderTimestamp time.Time
	windowNewerTimestamp time.Time
	nextQuery            url.Values
	resultSetOrder       ResultSetOrder
	pageSize             int
	client               *Client
	hasNextPage          bool
}

// Rewind rewinds the paginator to the first page. GetNextPage will then return the first page of the call.
func (d *DatastreamPaginator) Rewind() {
	// Invalid time
	d.windowOlderTimestamp = time.Time{}
	d.windowNewerTimestamp = time.Time{}
	d.nextQuery = url.Values{}
	d.hasNextPage = true
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
	//EHEH I WAS RIGHT: d.computePageState(len(page), page[len(page)-1].Timestamp)
	req := d.client.makeHTTPrequest(http.MethodGet, callURL, nil, d.client.token)

	return GetNextDatastreamPageRequest{req: req, expects: 200, paginator: d}, nil
}

type GetNextDatastreamPageRequest struct {
	req       *http.Request
	expects   int
	paginator Paginator
}

func (r GetNextDatastreamPageRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return GetNextDatastreamPageResponse{res: res, paginator: &r.paginator}, nil
}

func (r GetNextDatastreamPageRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

func (d *DatastreamPaginator) setupCallURL() (*url.URL, error) {
	// TODO check err
	callURL, _ := url.Parse(d.baseURL.String())

	query := d.nextQuery
	if d.resultSetOrder == AscendingOrder {
		if d.pageSize != 0 {
			return &url.URL{}, fmt.Errorf("A limit parameter must be specified when using AscendingOrder")
		}
		query.Set("limit", fmt.Sprintf("%d", d.pageSize))
		// check that a last value does actually exist before setting 'to'
		if (d.windowOlderTimestamp != time.Time{}) {
			// All data in the next page
			// come from a time until 'to' (so we ascend)
			query.Set("to", d.windowOlderTimestamp.UTC().Format(time.RFC3339Nano))
		}
	} else {
		// If no start is set, let's start from the beginnning of time
		if (d.windowOlderTimestamp == time.Time{}) {
			d.windowOlderTimestamp = time.Unix(0, 0)
		}
		// All data in the next page
		// come from a time after 'since' (so we descend)
		query.Set("since", d.windowOlderTimestamp.UTC().Format(time.RFC3339Nano))
		if (d.windowNewerTimestamp != time.Time{}) {
			// All data in the next page
			// come from a time until 'to'
			query.Set("to", d.windowNewerTimestamp.UTC().Format(time.RFC3339Nano))
		}
		if d.pageSize != 0 {
			query.Set("limit", fmt.Sprintf("%d", d.pageSize))
		}

	}
	callURL.RawQuery = query.Encode()

	return callURL, nil
}
