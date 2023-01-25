package newclient

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"moul.io/http2curl"
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

type GetNextDeviceListPageRequest struct {
	req       *http.Request
	expects   int
	paginator Paginator
}

// Performs a request to get the next page.
// Returns either a response that can be parsed with Parse() or an error if the request failed.
func (r GetNextDeviceListPageRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return GetNextDeviceListPageResponse{res: res, paginator: &r.paginator}, nil
}

// Returns the curl command corresponding to the request to get the next page.
func (r GetNextDeviceListPageRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

// GetNextPage returns a request to get the next result page from the paginator.
// If no more results are available, HasNextPage will return false.
// GetNextPage throws an error if no more pages are available.
func (d *DeviceListPaginator) GetNextPage() (AstarteRequest, error) {
	if !d.hasNextPage {
		return Empty{}, errors.New("No more pages available")
	}

	callURL := d.setupCallURL()
	req := d.client.makeHTTPrequest(http.MethodGet, callURL, nil, d.client.token)

	return GetNextDeviceListPageRequest{req: req, expects: 200, paginator: d}, nil
}

func (d *DeviceListPaginator) setupCallURL() *url.URL {
	// TODO check err
	callURL, _ := url.Parse(d.baseURL.String())
	query := d.nextQuery
	switch d.format {
	case DeviceIDFormat:
		query.Set("details", "false")
	case DeviceDetailsFormat:
		query.Set("details", "true")
	}

	callURL.RawQuery = query.Encode()

	return callURL
}
