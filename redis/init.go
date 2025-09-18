package redis

import (
	"fmt"

	goRedis "github.com/go-redis/redis"
	"github.com/guestin/kboot"
	"github.com/guestin/log"
)

var logger log.ClassicLog
var zapLogger log.ZapLog

func init() {
	kboot.RegisterUnit("redis", _init)
}

func _init(unit kboot.Unit) (kboot.ExecFunc, error) {
	logger = unit.GetClassicLogger()
	zapLogger = unit.GetZapLogger()
	_cfg = new(Config)
	err := unit.UnmarshalSubConfig(ModuleName, _cfg,
		kboot.MustBindEnv(CfgKeyHost),
		kboot.MustBindEnv(CfgKeyPort),
		kboot.MustBindEnv(CfgKeyDbIdx),
		kboot.MustBindEnv(CfgKeyPassword),
		kboot.MustBindEnv(CfgKeyKeyPrefix),
	)
	if err != nil {
		return nil, err
	}
	if _cfg.Port == 0 {
		_cfg.Port = 6379
	}
	if _cfg.KeyPrefix == "" {
		_cfg.KeyPrefix = DefaultKeyPrefix
	}
	option := &goRedis.Options{
		Addr:     fmt.Sprintf("%s:%d", _cfg.Host, _cfg.Port),
		Password: _cfg.Password,
		DB:       _cfg.DbIdx,
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
