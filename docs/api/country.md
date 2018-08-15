# Country Methods (API)

The Country API provides an interface to query information about the individual country models.

## List Available Country Models

----

Get a list of available country models

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/countries | 200         | application/json |

### Parameters

None

### Sample Request

```
$ curl -X GET "http://{HOST}:{PORT}/v1/countries"
```

### Sample Response

```json
[
  "US"
]
```

## Notes

Detailed information about the attributes and heuristics associated with the model can be accessing the [attribute APIs](attribute.md) and [heuristics APIs](heuristic.md) within the countries path.
