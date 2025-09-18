package db

import (
	"fmt"

	"github.com/guestin/kboot"
	"github.com/guestin/log"
	"github.com/guestin/mob/merrors"
)

var logger log.ClassicLog
var zapLogger log.ZapLog

func init() {
	kboot.RegisterUnit(ModuleName, _init)
}

func _init(unit kboot.Unit) (kboot.ExecFunc, error) {
	logger = unit.GetClassicLogger()
	zapLogger = unit.GetZapLogger()
	cfgList, err := bindConfig(unit)
	if err != nil {
		return nil, err
	}
	if len(cfgList) == 0 {
		return nil, merrors.Errorf("no valid db Config found")
	}
	for _, cfg := range cfgList {
		ds := cfg.name
		orm, err := newORM(unit.GetContext(), *cfg)
		if err != nil {
			return nil, merrors.Errorf("init datasource '%s' err : %v", ds, err)
		}
		_ormMaps.Store(ds, orm)
		if ds == cfgKeyDefault {
			_ormDB = orm
		}
	}

	if _migrator != nil {
		err := _migrator()
		if err != nil {
			return nil, merrors.Errorf("migrate error : %v", err)
		}
	}
	return _execute, nil
}

func bindConfig(unit kboot.Unit) (map[string]*Config, error) {
	ret := make(map[string]*Config)
	defaultCfg := new(Config)
	err := unit.UnmarshalSubConfig(ModuleName, defaultCfg,
		kboot.MustBindEnv(cfgKeyDbDsn),
		kboot.MustBindEnv(cfgKeyDbDebug),
		kboot.MustBindEnv(cfgKeyDbType),
		kboot.MustBindEnv(cfgKeyDbTimezone),
	)
	if err != nil {
		return nil, err
	}
	defaultCfg.name = cfgKeyDefault
	ret[cfgKeyDefault] = defaultCfg
	dbSettings := unit.GetRawViper().Sub(ModuleName).AllSettings()
	for key, _ := range dbSettings {
		switch dbSettings[key].(type) {
		case map[string]interface{}:
			_, exist := ret[key]
			if exist {
				return nil, merrors.Errorf("duplicate db setting : %s", key)
			}
			logger.Infof("try parser db settings '%s' ...", key)
			var dsCfg = new(Config)
			if err = unit.UnmarshalSubConfig(fmt.Sprintf("%s.%s", ModuleName, key), dsCfg,
				kboot.MustBindEnv(cfgKeyDbDsn),
				kboot.MustBindEnv(cfgKeyDbDebug),
				kboot.MustBindEnv(cfgKeyDbType),
				kboot.MustBindEnv(cfgKeyDbTimezone),
			); err != nil {
				return nil, err
			}
			dsCfg.name = key
			ret[key] = dsCfg
		}
	}
	return ret, nil
}

func _execute(unit kboot.Unit) kboot.ExitResult {
	<-unit.Done()
	return kboot.ExitResult{
		Code:  0,
		Error: nil,
	}
}
