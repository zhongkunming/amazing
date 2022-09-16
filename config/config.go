package config

type Config struct {
	App App `yaml:"app" json:"app"`

	DailyCheckConfig DailyCheckConfig `yaml:"daily_check" json:"daily_check"`
}
