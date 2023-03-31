package config

type ServerConfig struct {
	Name              string            `mapstructure:"name" json:"name"`
	Port              int               `mapstructure:"port" json:"port"`
	Tags              []string          `mapstructure:"tags" json:"tags"`
	DBConfig          DBConfig          `mapstructure:"db" json:"db"`
	ESConfig          ESConfig          `mapstructure:"es" json:"es"`
	RedisConfig       RedisConfig       `mapstructure:"redis" json:"redis"`
	JWTConfig         JWTConfig         `mapstructure:"jwt" json:"jwt"`
	OauthQQConfig     OauthQQConfig     `mapstructure:"oauth_qq_config" json:"oauth_qq_config"`
	OauthGithubConfig OauthGithubConfig `mapstructure:"oauth_github_config" json:"oauth_github_config"`
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

type JWTConfig struct {
	Issuer     string `mapstructure:"issuer" json:"issuer"`
	Expire     int64  `mapstructure:"expire" json:"expire"`
	SigningKey string `mapstructure:"key" json:"signingKey"`
}

type OauthQQConfig struct {
	AppId     string `mapstructure:"app_id" json:"app_id"`
	AppSecret string `mapstructure:"app_secret" json:"app_secret"`
	Uri       string `mapstructure:"uri" json:"uri"`
}

type OauthGithubConfig struct {
	AppId     string `mapstructure:"app_id" json:"app_id"`
	AppSecret string `mapstructure:"app_secret" json:"app_secret"`
	Uri       string `mapstructure:"uri" json:"uri"`
}
