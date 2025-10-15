package kboot

import "github.com/spf13/viper"

type CfgOption Option[*viper.Viper]

func MustBindEnv(input ...string) CfgOption {
	return optionFunc[*viper.Viper](func(v *viper.Viper) {
		v.MustBindEnv(input...)
	})
}
