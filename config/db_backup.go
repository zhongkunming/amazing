package config

type DbBackup struct {
	BDb DbBackupDb `yaml:"bDb" json:"bDb"`
	SDb DbBackupDb `yaml:"sDb" json:"sDb"`
}

type DbBackupDb struct {
	Host     string `yaml:"host" json:"host"`
	Username string `yaml:"username" json:"username"`
	Passwd   string `yaml:"passwd" json:"passwd"`
	Database string `yaml:"database" json:"database"`
}
