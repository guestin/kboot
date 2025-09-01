package redis

import (
	"fmt"

	goRedis "github.com/go-redis/redis"
	"github.com/guestin/kboot"
	"github.com/guestin/log"
	"github.com/guestin/mob/merrors"
)

var logger log.ClassicLog
var zapLogger log.ZapLog

func init() {
	kboot.RegisterUnit("redis", _init)
}

func _init(unit kboot.Unit) (kboot.ExecFunc, error) {
	logger = unit.GetClassicLogger()
	zapLogger = unit.GetZapLogger()
	v := unit.GetGlobalContext().GetViper()
	host := v.GetString(buildCfgKey(CfgKeyHost))
	password := v.GetString(buildCfgKey(CfgKeyPassword))
	if host == "" {
		return nil, merrors.Errorf("redis.host is requierd")
	}
	port := 6379
	if v.IsSet(buildCfgKey(CfgKeyPort)) {
		port = v.GetInt(buildCfgKey(CfgKeyPort))
	}
	if port <= 0 || port > 65535 {
		return nil, merrors.Errorf("redis.port is not valid , must be (0,65535]")
	}
	db := 0
	if v.IsSet(buildCfgKey(CfgKeyDbIdx)) {
		db = v.GetInt(buildCfgKey(CfgKeyDbIdx))
	}
	if db < 0 {
		return nil, merrors.Errorf("redis.db is not valid , must >= 0")
	}
	redisAddr := fmt.Sprintf("%s:%d", host, port)
	option := &goRedis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	}
	_redisCli = &Client{Client: goRedis.NewClient(option)}
	return _execute, nil
}

func _execute(unit kboot.Unit) kboot.ExitResult {
	<-unit.Done()
	return kboot.ExitResult{
		Code:  0,
		Error: nil,
	}
}
