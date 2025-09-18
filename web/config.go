package web

import (
	"fmt"

	"github.com/guestin/kboot/web/mid"
)

const (
	ModuleName = "web"

	CfgKeyListen = "listen"
	CfgKeyDebug  = "debug"

	DefaultListenAddress = ":8080"
)

func buildCfgKey(key string) string {
	return fmt.Sprintf("%s.%s", ModuleName, key)
}

type (
	Config struct {
		ListenAddress string          `toml:"listen" validate:"required" mapstruct:"auth"`
		Debug         bool            `toml:"debug"`
		Auth          *mid.AuthConfig `toml:"auth" validate:"omitnil" mapstruct:"auth"`
	}
)
