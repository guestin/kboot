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
	v := unit.GetGlobalContext().GetViper()

	listenAddress := DefaultListenAddress
	if v.IsSet(buildCfgKey(CfgKeyListen)) {
		listenAddress = v.GetString(buildCfgKey(CfgKeyListen))
	}
	debug := false
	if v.IsSet(buildCfgKey(CfgKeyDebug)) {
		debug = v.GetBool(buildCfgKey(CfgKeyDebug))
	}
	return _initEcho(unit, config{
		ListenAddress: listenAddress,
		Debug:         debug,
	})
}
