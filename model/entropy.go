package model

type heuristic struct {
	Match  []string `json:"match"`
	Insert []string `json:"insert"`
	Remove []string `json:"remove"`
}

type heuristics []heuristic

type format struct {
	Format string `json:"format"`
	Score  int    `json:"score"`
}

type formats []format

type attribute struct {
	Name    string `json:"name"`
	Formats formats `json:"formats"`
}

type attributes []attribute

type Entropy struct {
	Locale     string  `json:"locale"`
	Threshold  float64 `json:"threshold"`
	K          int     `json:"k"`
	Heuristics heuristics `json:"heuristics"`
	Attributes attributes `json:"attributes"`
}
