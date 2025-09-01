package db

import (
	"context"

	"github.com/guestin/kboot"
	"github.com/guestin/log"
	"github.com/guestin/mob/merrors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
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
		return nil, merrors.Errorf("no valid db config found")
	}
	//_, defaultExist := cfgList[CfgKeyDefault]
	//if !defaultExist {
	//	return nil, merrors.Errorf("default db not configured")
	//}
	for _, cfg := range cfgList {
		ds := cfg.name
		orm, err := newORM(unit.GetContext(), *cfg)
		if err != nil {
			return nil, merrors.Errorf("init datasource '%s' err : %v", ds, err)
		}
		_ormMaps.Store(ds, orm)
		if ds == CfgKeyDefault {
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

func bindConfig(unit kboot.Unit) (map[string]*config, error) {
	v := unit.GetGlobalContext().GetViper()
	ret := make(map[string]*config)
	defaultCfg := new(config)
	dbV := v.Sub(ModuleName)
	if dbV == nil {
		return nil, nil
	}
	dbV.MustBindEnv(CfgKeyDbDsn)
	dbV.MustBindEnv(CfgKeyDbDebug)
	dbV.MustBindEnv(CfgKeyDbType)
	dbV.MustBindEnv(CfgKeyDbTimezone)
	var err error
	logger.Infof("try parser default db settings ...")
	err = dbV.Unmarshal(defaultCfg)
	if err != nil {
		return nil, err
	}
	if err = defaultCfg.validate(); err == nil {
		defaultCfg.name = CfgKeyDefault
		ret[CfgKeyDefault] = defaultCfg
	}
	dbSettings := dbV.AllSettings()
	for key, _ := range dbSettings {
		switch dbSettings[key].(type) {
		case map[string]interface{}:
			_, exist := ret[key]
			if exist {
				return nil, merrors.Errorf("duplicate db setting : %s", key)
			}
			logger.Infof("try parser db settings '%s' ...", key)
			dsV := dbV.Sub(key)
			dsV.MustBindEnv(CfgKeyDbDsn)
			dsV.MustBindEnv(CfgKeyDbDebug)
			dsV.MustBindEnv(CfgKeyDbType)
			dsV.MustBindEnv(CfgKeyDbTimezone)
			var dsCfg = new(config)
			if err = dsV.Unmarshal(dsCfg); err != nil {
				return nil, merrors.Errorf("invalid db setting : %s , %v", key, err)
			}
			if err = dsCfg.validate(); err != nil {
				return nil, merrors.Errorf("invalid db setting : %s , %v", key, err)
			}
			dsCfg.name = key
			ret[key] = dsCfg
		}
	}
	return ret, nil
}

func buildOrm(ctx context.Context, v *viper.Viper, ds string) (*gorm.DB, error) {
	logger.Infof("init orm '%s'", ds)
	dbType := v.GetString(buildCfgKey(CfgKeyDbType, ds))
	dbDsn := v.GetString(buildCfgKey(CfgKeyDbDsn, ds))
	dbDebug := v.GetBool(buildCfgKey(CfgKeyDbDebug, ds))
	dbTz := v.GetString(buildCfgKey(CfgKeyDbTimezone, ds))
	if dbType == "" {
		dbType = DsTypePg
	}
	if dbDsn == "" {
		return nil, merrors.Errorf("db.dsn is required for datasource '%s'", ds)
	}
	cfg := config{
		Type:     dbType,
		DSN:      dbDsn,
		Debug:    dbDebug,
		Timezone: dbTz,
	}
	orm, err := newORM(ctx, cfg)
	if err != nil {
		return nil, merrors.Errorf("init datasource '%s' err : %v", ds, err)
	}
	return orm, nil
}

func _execute(unit kboot.Unit) kboot.ExitResult {
	<-unit.Done()
	return kboot.ExitResult{
		Code:  0,
		Error: nil,
	}
}
