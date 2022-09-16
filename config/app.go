package config

type App struct {
	Name        string `yaml:"name" json:"name"`
	LogFileName string `yaml:"logFileName" json:"logFileName"`
}
