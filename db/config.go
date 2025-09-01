package db

import (
	"fmt"

	"github.com/guestin/mob/merrors"
)

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

type config struct {
	name     string
	Type     string `toml:"type" mapstructure:"type"`
	DSN      string `toml:"dsn" mapstructure:"dsn"`
	Debug    bool   `toml:"debug" mapstructure:"debug"`
	Timezone string `toml:"timezone" mapstructure:"timezone"`
}

func (this *config) validate() error {
	if this.Type != "" && this.Type != DsTypePg && this.Type != DsTypeSqlLite {
		return merrors.Errorf("invalid db type : %s , must be %s or %s", this.Type, DsTypePg, DsTypeSqlLite)
	}
	if this.DSN == "" {
		return merrors.Errorf("dsn must not be empty")
	}
	return nil
}

func buildCfgKey(key string, dsName ...string) string {
	if len(dsName) > 0 && dsName[0] != "" {
		return fmt.Sprintf("%s.%s.%s", ModuleName, dsName[0], key)
	}
	return fmt.Sprintf("%s.%s", ModuleName, key)
}
