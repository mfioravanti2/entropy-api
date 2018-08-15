# Metric Methods (API)

The Entropy APIs provide a mechanism for collecting and querying system status.

## Get System Health

----

Get details about the current system status.

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/sys/health | 200         | application/json |

### Parameters

None

### Sample Request

```
$ curl -X GET "http://{HOST}:{PORT}/v1/sys/health"
```

### Sample Response

```json
{
    "status": "good",
    "api_version": "0.0.1",
    "model_versions": [
        {
            "country": "US",
            "timestamp": "2018-06-19T22:32:31.366105009Z",
            "version": "0.0.18"
        }
    ],
    "data_store": {
        "status": "good",
        "engine": "none",
        "last_use": "2018-08-15T13:47:21.924442409Z"
    }
}
```

## Get OpenAPI v3.0 Specification

----

Get the OpenAPI v3 (Swagger) Documentation for the Entropy API.

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/sys/spec  | 200         | application/json |

### Parameters

None

### Sample Request

```
$ curl -X GET "http://{HOST}:{PORT}/v1/sys/spec"
```

## Get System Metrics

Metrics about the endpoints are their usage is available through the Metrics endpoint.  More information about that endpoint is available in the [Metrics API documentation](metrics).

## Reload Models

----

Reload the country models

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/sys/reload | 200         |  |

### Parameters

None

### Sample Request

```
$ curl -X GET "http://{HOST}:{PORT}/v1/sys/reload"
```
