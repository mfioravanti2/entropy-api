# Scoring Methods (API)

The Entropy APIs provide a mechanism for collecting usage metrics.

## Score an Attribute Set

----

Get a list of available endpoint metrics

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| POST   | /v1/scores    | 200         | application/json |

### Query Parameters

* **mode** (string:[**detailed**|summary] <optional>) - Specify if the response should contain detailed information about the scoring process or only summary level information. The detailed calculation method is the default.
* **format** (string:[naive|**mean**|rare] <optional>) - Specify the method of calculation for scoring all attributes regardless of the format specified in the scoring request. The default scoring format is mean.
* **reductions** (string:[**include**|exclude] <optional>) - Specify if the reduction heuristics should be included in the scoring process. The default scoring process includes applying the heuristics to reduce the duplicate information.


### Sample Payload

```json
{
  "locale" : "US",
  "people" : [
    {
      "nationality" : "US",
      "person_id" : "0",
      "attributes" : [
        { "mnemonic" : "phone.nanpa.full", "format" : "mean", "tag" : "work" },
        { "mnemonic" : "phone.nanpa.full", "format" : "mean", "tag" : "personal" }
      ]
    }
  ]
}
```

### Sample Request

```
$ curl -X POST --data @payload.json "http://{HOST}:{PORT}/v1/scores"
```

### Sample Response

```json
{
    "data": {
        "pii": true,
        "locale": "US",
        "score": 47.6232,
        "api_version": "0.0.1",
        "run_date": "2018-08-15T18:10:24.114536348Z",
        "people": [
            {
                "id": "0",
                "nationality": "US",
                "score": 47.6232,
                "attributes": [
                    {
                        "mnemonic": "phone.nanpa.country_code",
                        "tag": "work",
                        "locale": "US",
                        "format": "mean",
                        "score": 0
                    },
                    {
                        "mnemonic": "phone.nanpa.area_code",
                        "tag": "work",
                        "locale": "US",
                        "format": "mean",
                        "score": 4.2461
                    },
                    {
                        "mnemonic": "phone.nanpa.central_office",
                        "tag": "work",
                        "locale": "US",
                        "format": "mean",
                        "score": 9.5997
                    },
                    {
                        "mnemonic": "phone.nanpa.subscriber_number",
                        "tag": "work",
                        "locale": "US",
                        "format": "mean",
                        "score": 9.9658
                    },
                    {
                        "mnemonic": "phone.nanpa.country_code",
                        "tag": "personal",
                        "locale": "US",
                        "format": "mean",
                        "score": 0
                    },
                    {
                        "mnemonic": "phone.nanpa.area_code",
                        "tag": "personal",
                        "locale": "US",
                        "format": "mean",
                        "score": 4.2461
                    },
                    {
                        "mnemonic": "phone.nanpa.central_office",
                        "tag": "personal",
                        "locale": "US",
                        "format": "mean",
                        "score": 9.5997
                    },
                    {
                        "mnemonic": "phone.nanpa.subscriber_number",
                        "tag": "personal",
                        "locale": "US",
                        "format": "mean",
                        "score": 9.9658
                    }
                ],
                "heuristics": [
                    "02294b1a-4bb7-443e-8b6c-768f303769da",
                    "02294b1a-4bb7-443e-8b6c-768f303769da"
                ]
            }
        ]
    }
}
```
