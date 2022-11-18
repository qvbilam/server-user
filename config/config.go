package config

type ServerConfig struct {
	Name        string      `mapstructure:"name" json:"name"`
	Port        int         `mapstructure:"port" json:"port"`
	Tags        []string    `mapstructure:"tags" json:"tags"`
	DBConfig    DBConfig    `mapstructure:"db" json:"db"`
	ESConfig    ESConfig    `mapstructure:"es" json:"es"`
	RedisConfig RedisConfig `mapstructure:"redis" json:"redis"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Database string `mapstructure:"database" json:"database"`
}

type ESConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Database int    `mapstructure:"database" json:"database"`
}
