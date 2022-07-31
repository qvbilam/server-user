package config

type ServerConfig struct {
	Name     string   `mapstructure:"name" json:"name"`
	Tags     []string `mapstructure:"tags" json:"tags"`
	DBConfig DBConfig `mapstructure:"db" json:"db"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Database string `mapstructure:"database" json:"database"`
}
