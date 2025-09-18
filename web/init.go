package web

import (
	"github.com/guestin/kboot"
	"github.com/guestin/kboot/web/internal"
	"github.com/guestin/log"
)

var logger log.ClassicLog
var zapLogger log.ZapLog

func init() {
	kboot.RegisterUnit(ModuleName, _init)
}

func _init(unit kboot.Unit) (kboot.ExecFunc, error) {
	internal.Log = unit.GetClassicLogger()
	internal.ZapLog = unit.GetZapLogger()
	logger = unit.GetClassicLogger()
	zapLogger = unit.GetZapLogger()
	cfg := new(Config)
	err := unit.UnmarshalSubConfig(ModuleName, cfg,
		kboot.MustBindEnv(CfgKeyListen),
		kboot.MustBindEnv(CfgKeyDebug),
	)
	if err != nil {
		return nil, err
	}
	return _initEcho(unit, cfg)
}
