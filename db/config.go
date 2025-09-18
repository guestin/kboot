package db

const (
	ModuleName = "db"

	CfgKeyDefault    = "default"
	CfgKeyDbType     = "type"
	CfgKeyDbDsn      = "dsn"
	CfgKeyDbDebug    = "debug"
	CfgKeyDbTimezone = "timezone"

	DsTypePg      = "postgres"
	DsTypeSqlLite = "sqlite"
)

type Config struct {
	name     string
	Type     string `toml:"type" validate:"required,oneof=postgres sqlite" mapstructure:"type"`
	DSN      string `toml:"dsn" validate:"required" mapstructure:"dsn"`
	Debug    bool   `toml:"debug" mapstructure:"debug"`
	Timezone string `toml:"timezone" mapstructure:"timezone"`
}
