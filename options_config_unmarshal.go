package kboot

import "github.com/spf13/viper"

type CfgUnmarshalOption Option[*viper.Viper]

func MustBindEnv(input ...string) CfgUnmarshalOption {
	return optionFunc[*viper.Viper](func(v *viper.Viper) {
		v.MustBindEnv(input...)
	})
}
