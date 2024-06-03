package model

type Config struct {
	Name    string            `json:"name"`
	Version int               `json:"version"`
	Params  map[string]string `json:"params"`
}

type ConfigRepository interface {
	Add(config Config) error
	Get(name string, version int) (Config, error)
	Delete(name string, version int) error
	NewConfigFromLiteral(literal string) (Config, error)
}
