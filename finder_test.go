package kboot

import (
	"testing"

	"github.com/guestin/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestAutoFindConfig(t *testing.T) {
	rootLogger, _ := log.EasyInitConsoleLogger(zap.DebugLevel, zap.DPanicLevel)
	logger := log.NewTaggedZapLogger(rootLogger, "test")
	fs := findConfigFiles(logger, []string{"./config"}, "", viper.SupportedExts)
	t.Log("fs", fs)
}
