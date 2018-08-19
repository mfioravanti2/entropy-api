package config

type Listener struct {
	Host string				`json:"host"`
	Port int				`json:"port"`
	UseTLS bool				`json:"use_TLS"`
	Encryption *Encryption	`json:"encryption,omitempty"`
}

