# Metric Methods (API)

The Entropy APIs provide a mechanism for collecting usage metrics.

## Get Endpoint Metrics

----

Get a list of available endpoint metrics

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/sys/metrics | 200         | application/json |

### Parameters

None

### Sample Request

```
$ curl -X GET "http://{HOST}:{PORT}/v1/sys/metrics"
```

### Sample Response

```json
{
    "entropy.attributes_details.get": {
        "count": 1
    },
    "entropy.attributes_details.get.status.200": {
        "count": 1
    },
    "entropy.attributes_list.get": {
        "count": 1
    },
    "entropy.attributes_list.get.status.200": {
        "count": 1
    },
    "entropy.country.get": {
        "count": 1
    },
    "entropy.country.get.status.200": {
        "count": 1
    },
    "entropy.heuristics_details.get": {
        "count": 1
    },
    "entropy.heuristics_details.get.status.200": {
        "count": 1
    },
    "entropy.heuristics_list.get": {
        "count": 1
    },
    "entropy.heuristics_list.get.status.200": {
        "count": 1
    },
    "entropy.sys.metrics.get": {
        "count": 1
    }
}
```
