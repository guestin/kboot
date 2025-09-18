package redis

const (
	ModuleName       = "redis"
	DefaultPort      = 6379
	DefaultKeyPrefix = "cn.guestin.kboot"

	cfgKeyHost      = "host"
	cfgKeyPort      = "port"
	cfgKeyDbIdx     = "db"
	cfgKeyPassword  = "password"
	cfgKeyKeyPrefix = "keyPrefix"
)

type Config struct {
	Host      string `toml:"host" validate:"required" mapstructure:"host"`
	Port      int    `toml:"port" validate:"omitempty,gte=0,lte=65535" mapstructure:"port"`
	DbIdx     int    `toml:"db" validate:"omitempty,gte=0,lte=128" mapstructure:"db"`
	Password  string `toml:"password" mapstructure:"password"`
	KeyPrefix string `toml:"keyPrefix" validate:"omitempty,excludesall= !@#$%^&*()\t\n\r" mapstructure:"keyPrefix"`
}

var _cfg *Config

func GetConfig() Config {
	return *_cfg
}
