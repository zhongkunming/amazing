package config

type Flows struct {
	LoginUrl string      `yaml:"loginUrl" json:"loginUrl"`
	FlowsUrl string      `yaml:"flowsnUrl" json:"flowsnUrl"`
	Users    []FLowsUser `yaml:"users" json:"users"`
}

type FLowsUser struct {
	Email      string `yaml:"email" json:"email"`
	Passwd     string `yaml:"passwd" json:"passwd"`
	RememberMe string `yaml:"remember_me" json:"remember_me" default:"on"`
	Code       string `yaml:"code" json:"code" default:""`
}
