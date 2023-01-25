package newclient

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/astarte-platform/astarte-go/interfaces"
	"moul.io/http2curl"
)

type GetDatastreamSnapshotRequest struct {
	req         *http.Request
	expects     int
	aggregation interfaces.AstarteInterfaceAggregation
}

const defaultPageSize = 100

// GetDatastreamSnapshot builds a request to return all the last values on all paths for a Datastream individual aggregate interface.
func (c *Client) GetDatastreamIndividualSnapshot(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	interfaceName string) (AstarteRequest, error) {

	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/%s", realm, devicePath(deviceIdentifier, deviceIdentifierType), interfaceName))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)

	return GetDatastreamSnapshotRequest{req: req, expects: 200, aggregation: interfaces.IndividualAggregation}, nil
}

// GetDatastreamAggregateSnapshot  builds a request to return the last value for a Datastream object aggregate interface
func (c *Client) GetDatastreamAggregateSnapshot(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	interfaceName string) (AstarteRequest, error) {

	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/%s", realm, devicePath(deviceIdentifier, deviceIdentifierType), interfaceName))

	// Quirk: Astarte returns all data, we must limit to the first one
	query := url.Values{}
	query.Set("limit", fmt.Sprintf("%d", 1))
	callURL.RawQuery = query.Encode()

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)

	return GetDatastreamSnapshotRequest{req: req, expects: 200, aggregation: interfaces.ObjectAggregation}, nil
}

func (r GetDatastreamSnapshotRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return GetDatastreamSnapshotResponse{res: res, aggregation: r.aggregation}, nil
}

func (r GetDatastreamSnapshotRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

// GetDatastreamsPaginator returns a Paginator for all the values on a path for a Datastream interface.
func (c *Client) GetDatastreamsPaginator(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, resultSetOrder ResultSetOrder) (Paginator, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	return c.getDatastreamPaginator(realm, deviceIdentifier, resolvedDeviceIdentifierType, interfaceName, interfacePath, time.Time{}, time.Now(), defaultPageSize, resultSetOrder)
}

// GetDatastreamsTimeWindowPaginator returns a Paginator for all the values on a path in a specified time window for a Datastream interface.
func (c *Client) GetDatastreamsTimeWindowPaginator(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, since, to time.Time, resultSetOrder ResultSetOrder) (Paginator, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	return c.getDatastreamPaginator(realm, deviceIdentifier, resolvedDeviceIdentifierType, interfaceName, interfacePath, since, to, defaultPageSize, resultSetOrder)
}

func (c *Client) getDatastreamPaginator(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string,
	since, to time.Time, pageSize int, resultSetOrder ResultSetOrder) (Paginator, error) {

	baseURL, _ := url.Parse(c.appEngineURL.String())
	baseURL.Path = path.Join(baseURL.Path, fmt.Sprintf("/v1/%s/%s/%s", realm, devicePath(deviceIdentifier, deviceIdentifierType), interfaceName))

	datastreamPaginator := DatastreamPaginator{
		baseURL:              baseURL,
		windowNewerTimestamp: time.Time{},
		windowOlderTimestamp: time.Time{},
		nextQuery:            url.Values{},
		pageSize:             pageSize,
		client:               c,
		hasNextPage:          true,
		resultSetOrder:       resultSetOrder,
	}

	if resultSetOrder == AscendingOrder {
		if pageSize != 0 {
			return &DatastreamPaginator{}, fmt.Errorf("A pageSize must be specified when using AscendingOrder")
		}
		if (since != time.Time{}) {
			return &DatastreamPaginator{}, fmt.Errorf("Specifying \"since\" is not supported when using AscendingOrder")
		}
		// check that a last value does actually exist before setting 'to'
		if (to != time.Time{}) {
			datastreamPaginator.windowOlderTimestamp = to
		}
	} else {
		// If no start is set, let's start from the beginnning of time
		if (since == time.Time{}) {
			datastreamPaginator.windowOlderTimestamp = time.Unix(0, 0)
		}
		// All data in the next page
		// come from a time after 'since' (so we descend)
		if (to != time.Time{}) {
			datastreamPaginator.windowNewerTimestamp = to

		}
	}

	return &datastreamPaginator, nil
}

type GetPropertiesRequest struct {
	req     *http.Request
	expects int
}

// GetAllProperties builds a request to return all the currently set Properties on a given interface.
func (c *Client) GetAllProperties(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	interfaceName string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/%s", realm, devicePath(deviceIdentifier, deviceIdentifierType), interfaceName))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)

	return GetPropertiesRequest{req: req, expects: 200}, nil
}

func (r GetPropertiesRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return GetPropertiesResponse{res: res}, nil
}

func (r GetPropertiesRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

// GetProperty builds a request to return the currently set Property on a given Interface at a given path.
func (c *Client) GetProperty(realm string, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	interfaceName string, interfacePath string) (AstarteRequest, error) {
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/%s/%s", realm, devicePath(deviceIdentifier, deviceIdentifierType), interfaceName, interfacePath))

	req := c.makeHTTPrequest(http.MethodGet, callURL, nil, c.token)

	return GetPropertiesRequest{req: req, expects: 200}, nil
}

// SendData builds a request to send data on the specified interface. It performs all validity checks on the Interface object before moving forward
// with the operation, as such it is assumed that the operation will be always validated on the client side. If you have access to a native
// Interface object, accessing this method rather than the lower level ones is advised.
// payload must match a compatible type for the Interface path. In case of an aggregate interface, payload *must* be a
// map[string]interface{}, and each payload will be individually checked.
func (c *Client) SendData(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType,
	astarteInterface interfaces.AstarteInterface, interfacePath string, payload any) (AstarteRequest, error) {
	// Perform a set of checks depending on the interface structure
	switch {
	case astarteInterface.Ownership == interfaces.DeviceOwnership:
		return Empty{}, fmt.Errorf("cannot send data to device-owned interface %s %d.%d", astarteInterface.Name, astarteInterface.MajorVersion, astarteInterface.MinorVersion)
	case astarteInterface.Type == interfaces.PropertiesType, astarteInterface.Aggregation == interfaces.IndividualAggregation:
		// In this case, validate the individual message
		if err := interfaces.ValidateIndividualMessage(astarteInterface, interfacePath, payload); err != nil {
			return Empty{}, err
		}
	case astarteInterface.Aggregation == interfaces.ObjectAggregation:
		aggregatePayload, ok := payload.(map[string]interface{})
		if !ok {
			return Empty{}, fmt.Errorf("Data sent to interfaces with object aggregation must be a map[string]interface{}")
		}
		if err := interfaces.ValidateAggregateMessage(astarteInterface, interfacePath, aggregatePayload); err != nil {
			return Empty{}, err
		}
	}

	// If we got here, it's time to do the right thing.
	switch {
	case astarteInterface.Type == interfaces.PropertiesType:
		return c.SetProperty(realm, deviceIdentifier, deviceIdentifierType, astarteInterface.Name, interfacePath, payload)
	case astarteInterface.Aggregation == interfaces.IndividualAggregation:
		return c.SendDatastream(realm, deviceIdentifier, deviceIdentifierType, astarteInterface.Name, interfacePath, payload)
	case astarteInterface.Aggregation == interfaces.ObjectAggregation:
		p, ok := payload.(map[string]any)
		if !ok {
			return Empty{}, fmt.Errorf("Invalid payload type for object-aggregated interface: payload must be a map, got %T", p)
		}
		return c.SendDatastream(realm, deviceIdentifier, deviceIdentifierType, astarteInterface.Name, interfacePath, p)
	}

	// We should never get here
	return Empty{}, fmt.Errorf("Interface %s %d.%d has malformed type or aggregation", astarteInterface.Name, astarteInterface.MajorVersion, astarteInterface.MinorVersion)
}

type SendDatastreamRequest struct {
	req     *http.Request
	expects int
}

// SendDatastream builds a request to send a datastream to the given interface without additional checks.
// payload must be of a type compatible with the interface's endpoint. Any errors will be returned on the server side or
// in payload marshaling. If you have a native AstarteInterface object, calling SendData is advised
func (c *Client) SendDatastream(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, payload any) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/interfaces/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType), interfaceName+interfacePath))

	normalizedPayload := interfaces.NormalizePayload(payload, true)
	body, _ := makeBody(normalizedPayload)
	req := c.makeHTTPrequest(http.MethodPost, callURL, body, c.token)

	return SendDatastreamRequest{req: req, expects: 200}, nil
}

func (r SendDatastreamRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r SendDatastreamRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type SetPropertyRequest struct {
	req     *http.Request
	expects int
}

// SetProperty builds a request to set a property on the given interface without additional checks. payload must be of a type
// compatible with the interface's endpoint. Any errors will be returned on the server side or
// in payload marshaling. If you have a native AstarteInterface object, calling SendData is advised
func (c *Client) SetProperty(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName, interfacePath string, payload any) (AstarteRequest, error) {
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/interfaces/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType), interfaceName+interfacePath))

	normalizedPayload := interfaces.NormalizePayload(payload, true)
	body, _ := makeBody(normalizedPayload)
	req := c.makeHTTPrequest(http.MethodPut, callURL, body, c.token)

	return SetPropertyRequest{req: req, expects: 200}, nil
}

func (r SetPropertyRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r SetPropertyRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}

type UnsetPropertyRequest struct {
	req     *http.Request
	expects int
}

// UnsetProperty builds a request to delete a property on the given interface without additional checks.
func (c *Client) UnsetProperty(realm, deviceIdentifier string, deviceIdentifierType DeviceIdentifierType, interfaceName string, interfacePath string) (AstarteRequest, error) {
	// TODO check if mapping is unsettable
	resolvedDeviceIdentifierType := resolveDeviceIdentifierType(deviceIdentifier, deviceIdentifierType)
	callURL, _ := url.Parse(c.appEngineURL.String())
	callURL.Path = path.Join(callURL.Path, fmt.Sprintf("/v1/%s/%s/interfaces/%s", realm, devicePath(deviceIdentifier, resolvedDeviceIdentifierType), interfaceName+interfacePath))

	req := c.makeHTTPrequest(http.MethodDelete, callURL, nil, c.token)

	return UnsetPropertyRequest{req: req, expects: 204}, nil
}

func (r UnsetPropertyRequest) Run(c *Client) (AstarteResponse, error) {
	res, err := c.httpClient.Do(r.req)
	if err != nil {
		return Empty{}, err
	}
	if res.StatusCode != r.expects {
		return Empty{}, ErrDifferentStatusCode
	}
	return NoDataResponse{res: res}, nil
}

func (r UnsetPropertyRequest) ToCurl(c *Client) string {
	command, _ := http2curl.GetCurlCommand(r.req)
	return fmt.Sprint(command)
}
