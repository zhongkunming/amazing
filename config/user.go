package config

type User struct {
	Email      string `yaml:"email" json:"email"`
	Passwd     string `yaml:"passwd" json:"passwd"`
	RememberMe string `yaml:"remember_me" json:"remember_me" default:"on"`
	Code       string `yaml:"code" json:"code" default:""`
}
