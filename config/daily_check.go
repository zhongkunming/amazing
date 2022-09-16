package config

type DailyCheckConfig struct {
	LoginUrl   string                  `yaml:"loginUrl" json:"loginUrl"`
	CheckInUrl string                  `yaml:"checkInUrl" json:"checkInUrl"`
	Users      []DailyCheckUsersConfig `yaml:"users" json:"users"`
}

type DailyCheckUsersConfig struct {
	Email      string `yaml:"email" json:"email"`
	Passwd     string `yaml:"passwd" json:"passwd"`
	RememberMe string `yaml:"remember_me" json:"remember_me" default:"on"`
	Code       string `yaml:"code" json:"code" default:""`
}
