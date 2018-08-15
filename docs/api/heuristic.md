# Heuristic Methods (API)

The heuristic APIs expose additional information about the country model used by the scoring system. Heuristics are used to modify the attribute set before they are scored, this allowed duplicate information to be removed from the attribute set and prevent the resulting score from being artificially inflated.

## List Heuristics within a Country Model

----

Get a list of all available heuristics within a country model

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/countries/:countryCode/heuristics | 200         | application/json |

### Parameters

* **countryCode** (string: <required>) - A 2-digit Country Code as defined in ISO 3166-1 alpha-2. The country code is case-insensitive.

### Sample Request

Retrieve a list of all heuristic Ids within the "US" country model.

```
$ curl -X GET "http://{HOST}:{PORT}/v1/countries/us/heuristics"
```

### Sample Response

A JSON array with all of the heuristic identifiers for the specific country is returned.

```json
[
    "ca50582b-45c0-4746-bae7-7c2845f19399",
    "cdd56943-cc2b-4575-9328-5521e79e62d2",
    "e841d3bb-de84-42bf-b92f-d5d58a3a9a8d",
    "eb202cd8-7618-4be8-a946-df675e819cb7",
    "fbf8bd2a-5b6c-4f47-a779-f8749a4e0fe8"
]
```

## Get Details about a Specific Attribute

----

Get a list of available heuristics for a specific country model.

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/countries/:countryCode/heuristics/:heuristicId | 200         | application/json |

### Parameters

* **countryCode** (string: <required>) - A 2-digit Country Code as defined in ISO 3166-1 alpha-2. The country code is case-insensitive.
* **heuristicId** (string<uuid>: <required>) - An attribute mnemonic which specifies the attribute to be retrieved. All Heuristic IDs are UUID formatted strings.

### Sample Request

Retrieve the details about the heuristic (heuristicId: "bb99443e-990c-4278-8291-cc991681e406") for the US country model.

```
$ curl -X GET "http://{HOST}:{PORT}/v1/countries/us/heuristics/bb99443e-990c-4278-8291-cc991681e406"
```

### Sample Response

A JSON object with details about the requested heuristic from the country model is returned.


```json
{
    "id": "bb99443e-990c-4278-8291-cc991681e406",
    "notes": "",
    "match": [
        "ssn.area_number",
        "ssn.group_number",
        "ssn.serial_number"
    ],
    "insert": [
        "ssn"
    ],
    "remove": [
        "ssn.area_number",
        "ssn.group_number",
        "ssn.serial_number"
    ]
}
```
In this response, the heuristic indicates that if the attribute set contains {ssn.area_number, ssn.group_number, ssn.serial_number}, add the {ssn} attributes and then remove the {ssn.area_number, ssn.group_number, ssn.serial_number} attributes from the set.

## Notes

Heuristics operate on sets of Attributes, details about retrieving that information can be found in the [Attribute APIs](attribute.md) documentation.