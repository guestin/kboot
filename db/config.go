package db

const (
	ModuleName = "db"

	cfgKeyDefault    = "default"
	cfgKeyDbType     = "type"
	cfgKeyDbDsn      = "dsn"
	cfgKeyDbDebug    = "debug"
	cfgKeyDbTimezone = "timezone"

	DsTypePg      = "postgres"
	DsTypeSqlLite = "sqlite"
)

type Config struct {
	name  string
	Type  string `toml:"type" validate:"required,oneof=postgres sqlite" mapstructure:"type"`
	DSN   string `toml:"dsn" validate:"required" mapstructure:"dsn"`
	Debug bool   `toml:"debug" mapstructure:"debug"`
}
