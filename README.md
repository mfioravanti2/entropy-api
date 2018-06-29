# Entropy API

----

A binary approach to classifying attributes as either Personally Identifiable Information (PII) or not is commonly used. Companies may adopt formal definitions which explicitly list specific attributes and if any of these special attributes are found in a set, the resulting set is classified as PII. For example, attributes such as a US Social Security Number (SSN) is classified as PII, while an individual's given name (i.e. first name) by itself is not considered to be PII. [Sweeney](https://scholar.google.com/scholar?hl=en&as_sdt=0%2C7&q=k+anonymity+author%3Asweeney&btnG=) highlighted that the combination of various attributes which in and of themselves would not be considered PII, can be combined to uniquely identify an individual. Sweeney's work highlighted that the attribute set (US Postal ZIP Code, Date of Birth, and Gender) formed a quasi-identifier and could uniquely identify ~87% of the United States population.

Lodha and Thomas defined the concept of [Probabilistic Anonymity](https://scholar.google.com/scholar?hl=en&q=probabilistic+anonymity+author%3Alodha+author%3Athomas&btnG=) where various attributes could have their information content measured, and as a result of these estimates the likelihood of the set being classified as a quasi-identifier against a universal population could be measured. This research builds upon both the k-Anonymity model was well as the probabilistic anonymity model to allow attribute sets to be classified in real-time and aid in making a determination if enough information is present in an attribute set that the set should be classified as PII.

With the adoption of various privacy frameworks in the international community, such as the European Union's General Data Protection Regulation (GDPR) understanding when an attribute set is PII is of critical importance.  The Entropy API is a research tool designed to assist in determining when enough information is present in an attribute set that it should be classified as Personally Identifiable Information.

The key features of the Entropy API are:

* **REST API**: JSON payloads describing the attribute sets are accepted and classified based on the amount of information present in the attributes. Besides scoring an attribute set, the API offers endpoints which can be queried for additional information about the various attributes of the model.

* **Logging**: Detailed logging is available which illustrates the scoring and classification of the attribute sets.


### Installation

Entropy API requires a
[working Go 1.10+ installation](http://golang.org/doc/install) and a
properly set `GOPATH`.

```
$ go get -u github.com/mfioravanti2/entropy-api
```

will download and build the Entropy API tool, installing it in
`$GOPATH/bin/entropy-api`.
