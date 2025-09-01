package web

import "fmt"

const (
	ModuleName = "web"

	CfgKeyListen = "listen"
	CfgKeyDebug  = "debug"

	DefaultListenAddress = ":8080"
)

func buildCfgKey(key string) string {
	return fmt.Sprintf("%s.%s", ModuleName, key)
}

type config struct {
	ListenAddress string `toml:"listen"`
	Debug         bool   `toml:"debug"`
}
