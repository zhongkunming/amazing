package config

type Config struct {
	App App `yaml:"app" json:"app"`

	DailyCheck DailyCheck `yaml:"dailyCheck" json:"dailyCheck"`

	DbBackup DbBackup `yaml:"dbBackup" json:"dbBackup"`
}
