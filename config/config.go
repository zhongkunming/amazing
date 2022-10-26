package config

type Config struct {
	App App `yaml:"app" json:"app"`

	Flows Flows `yaml:"flows" json:"flows"`
}
