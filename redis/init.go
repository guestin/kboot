package redis

import (
	"github.com/guestin/kboot"
	"github.com/guestin/log"
	"github.com/guestin/mob/merrors"
	goRedis "github.com/redis/go-redis/v9"
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
		kboot.MustBindEnv(cfgKeyHost),
		kboot.MustBindEnv(cfgKeyPort),
		kboot.MustBindEnv(cfgKeyDbIdx),
		kboot.MustBindEnv(cfgKeyPassword),
		kboot.MustBindEnv(cfgKeyKeyPrefix),
	)
	if err != nil {
		return nil, err
	}
	if _cfg.KeyPrefix == "" {
		_cfg.KeyPrefix = DefaultKeyPrefix
	}
	option := &goRedis.UniversalOptions{
		Addrs:    _cfg.Address,
		DB:       _cfg.DbIdx,
		Password: _cfg.Password,
	}
	cli := goRedis.NewUniversalClient(option)
	if err = cli.Ping(unit.GetContext()).Err(); err != nil {
		return nil, merrors.ErrorWrap(err, "redis connect failed")
	}
	_redisCli = &Client{UniversalClient: cli}
	return _execute, nil
}

func _execute(unit kboot.Unit) kboot.ExitResult {
	<-unit.Done()
	return kboot.ExitResult{
		Code:  0,
		Error: nil,
	}
}
