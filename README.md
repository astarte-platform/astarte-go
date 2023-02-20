# Astarte Go

This is the main dependency module for all Astarte Go applications and SDKs.

Astarte Go requires at least Go 1.18.

_________________________

## Migrating from 0.90.x

Following https://github.com/astarte-platform/astarte-go/issues/30, this library exposes a new API.
The most significant changes are:
- A completely reworked `client` package
- The `misc` package has been split in `auth`, `deviceid` and `astarteservices` packages

### Using the new `client` package

This package provides functions to interact with Astarte APIs.

### Client setup

The creation of a new API client has been made clearer using the functional options pattern.
As an example, consider the following code, written using the v0.90.4 version of astarte-go:
```go
    astarteAPIClient, err = client.NewClient(astarteURL, nil)
    if err != nil {
        fmt.Println(err)
    }

    if err := astarteAPIClient.SetTokenFromPrivateKeyFileWithTTL(privateKeyFile, 60); err != nil {
        fmt.Println(err)
    }
```

From > v0.90.4, the same setup is handled in this way:
```go
    astarteAPIClient, err = client.New(
        client.WithBaseURL(astarteURL),
        client.WithPrivateKey(privateKeyFile),
        client.WithExpiry(60),
    )
    if err != nil {
        fmt.Println(err)
    }
```

### Astarte API

Moreover, functions handling Astarte APIs have been revisited. In general, the interaction is divided in three steps:
1. Generate a request, i.e. a value of type `AstarteRequest`
2. Perform the request (calling the `Run()` method on the request), thus obtaining a value of type `AstarteResponse`
3. Parse the result (calling the `Parse()` method on the response)

Each step may fail with an error, which is strongly recommended to check.
This pattern gives more control to users on how to handle each interaction step.
`AstarteRequest`s also provide a `ToCurl()` method to emit a command-line command equivalent to the request.
`AstarteResponse`s also provide a `Raw()` method to handle the response in an ad-hoc way.


As an example, consider the following code, written using the v0.90.4 version of astarte-go:
```go
    deviceDetails, err := astarteAPIClient.AppEngine.GetDevice(realm, deviceID, deviceIdentifierType)
    if err != nil {
        fmt.Println(err)
    }
```

From > v0.90.4, the same call can be performed as:
```go
    deviceDetailsReq, err := astarteAPIClient.GetDeviceDetails(realm, deviceID, deviceIdentifierType)
    if err != nil {
        fmt.Println(err)
    }

    deviceDetailsRes, err := deviceDetailsReq.Run(astarteAPIClient)
    if err != nil {
        fmt.Println(err)
    }

    deviceDetails, err := deviceDetailsRes.Parse()
    if err != nil {
        fmt.Println(err)
    }
```

### Pagination

Finally, pagination has been handled correctly.
Astarte returns paginated data in two cases: device lists and device data.
The `Paginator` interface provides methods to query those data, such as `HasNextPage()`, `GetNextPage()`, `GetPageSize()`.
The following example shows how to use the `DatastreamPaginator` to print a paginated set of samples on a datastream individual interface, 
starting from the oldest, with a page size of 10 .
```go
    paginator, err := c.GetIndividualDatastreamsPaginator(realm, deviceID, client.AstarteDeviceID, 
        "org.astarte.genericsensors.Value", "/streamTest/value", client.DescendingOrder, 10)
    if err != nil {
        fmt.Println(err)
    }

    for paginator.HasNextPage() {
        getNextPageReq, err := paginator.GetNextPage()
        if err != nil {
            fmt.Println(err)
        }

        getNextPageRes, err := getNextPageReq.Run(c)
        if err != nil {
            fmt.Println(err)
        }

        rawNextPageData, err := getNextPageRes.Parse()
        if err != nil {
            fmt.Println(err)
        }
        nextPageData, _ := rawNextPageData.([]client.DatastreamIndividualValue)
        for _, v := range nextPageData {
            fmt.Printf("Value: %#v, Timestamp: %#v, Reception Timestamp: %#v\n", v.Value, v.Timestamp, v.ReceptionTimestamp)
        }
    }
```

## Using the new `auth`, `deviceid` and `astarteservices` packages

Just replace `misc` with the new packages, the context of which is pretty much self-explainatory.