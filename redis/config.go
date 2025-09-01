package redis

import "fmt"

const (
	ModuleName = "redis"

	CfgKeyHost     = "host"
	CfgKeyPort     = "port"
	CfgKeyDbIdx    = "db"
	CfgKeyPassword = "password"
)

func buildCfgKey(key string) string {
	return fmt.Sprintf("%s.%s", ModuleName, key)
}
