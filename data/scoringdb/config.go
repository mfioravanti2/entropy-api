package scoringdb

type Config struct {
	Engine		string	`json:"engine"`
	Connection	string	`json:"connection"`
}

func NewConfig() (*Config, error) {
	c := Config{ Engine: "sqlite3", Connection: "./scores.db" }
	return &c, nil
}
