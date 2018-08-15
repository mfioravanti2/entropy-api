# Attribute Methods (API)

The attribute APIs expose additional information about the country model used by the scoring system. Each attribute has been "scored" (i.e. the frequency-based information theoretic entropy calculation has been performed). As the attributes are specific to each country, attributes can only be accessed by specifying a country code in the URL path.

## List Attributes within a Country Model

----

Get a list of all available attributes within a country models

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/countries/:countryCode/attributes | 200         | application/json |

### Parameters

* **countryCode** (string: <required>) - A 2-digit Country Code as defined in ISO 3166-1 alpha-2. The country code is case-insensitive.

### Sample Request

Retrieve a list of all mnemonics for the attributes within the "US" country model.

```
$ curl -X GET "http://{HOST}:{PORT}/v1/countries/us/attributes"
```

### Sample Response

A JSON array with all of the attribute's mnemonics for the specific country is returned.

```json
[
    "date_of_birth",
    "date_of_birth.age",
    "date_of_birth.birthday",
    "date_of_birth.day",
    "date_of_birth.month",
    "date_of_birth.year",
    "sex",
    "ssn",
    "ssn.area_number",
    "ssn.group_number",
    "ssn.serial_number"
]
```


## Get Details about a Specific Attribute

----

Get a list of available attributes for a specific country model.

| Method | Path          | Status Code | Content-Type     |
|--------|---------------|-------------|------------------|
| GET    | /v1/countries/:countryCode/attributes/:attributeMnemonic | 200         | application/json |

### Parameters

* **countryCode** (string: <required>) - A 2-digit Country Code as defined in ISO 3166-1 alpha-2. The country code is case-insensitive.
* **attributeMnemonic** (string: <required>) - An attribute mnemonic which specifies the attribute to be retrieved.  The mnemonic is case-sensitive.

### Sample Request

Retrieve the details about the Social Security Number (attributeMnemonic: ssn) for the US country model.

```
$ curl -X GET "http://{HOST}:{PORT}/v1/countries/us/attributes/ssn"
```

### Sample Response

A JSON object with details about the requested attribute from the country model is returned.


```json
{
    "id": "a4050423-c9ce-43ec-855c-cae9c76251ba",
    "name": "SSN",
    "mnemonic": "ssn",
    "notes": "",
    "sources": [
        {
            "title": "Social Security Numbers: The SSN Numbering Scheme",
            "organization": "U.S. Social Security Administration",
            "date": "2018-06-19T22:24:14Z",
            "url": "https://www.ssa.gov/history/ssn/geocard.html"
        },
        {
            "title": "Social Security Administration Fact Sheet: Issuing SSNs",
            "organization": "U.S. Social Security Administration",
            "date": "2011-06-13T16:31:56Z",
            "url": "https://www.ssa.gov/kc/SSAFactSheet--IssuingSSNs.pdf"
        },
        {
            "title": "IRS Publication 1346 (Rev. 10-2012): Electronic Return File Specifications and Record Layouts for Individual Tax Returns (Tax Year 2012)",
            "organization": "U.S. Internal Revenue Service",
            "date": "2012-10-16T15:15:22Z",
            "url": "https://www.irs.gov/pub/irs-pdf/p1346.pdf"
        }
    ],
    "formats": [
        {
            "format": "naive",
            "score": 29.7275
        },
        {
            "format": "mean",
            "score": 29.7275
        },
        {
            "format": "rare",
            "score": 29.7275
        }
    ]
}
```
