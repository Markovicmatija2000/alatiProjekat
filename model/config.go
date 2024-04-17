package model

type Config struct {
	Name    string            `json:"name"`
	Version int               `json:"version"`
	Params  map[string]string `json:"params"`
}

// TODO: Dodati metode

type ConfigRepository interface {
	Add(config Config)
	Get(name string, version int) (Config, error)
}
