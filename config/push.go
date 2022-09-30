package config

type Push struct {
	AppId string `yaml:"appId" json:"appId"`

	AppSecret string `yaml:"appSecret" json:"appSecret"`

	AccessToken string
}
