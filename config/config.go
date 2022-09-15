package config

type Config struct {
	App   App    `yaml:"app" json:"app"`
	Users []User `yaml:"users" json:"users"`
}
