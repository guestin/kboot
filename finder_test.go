package kboot

import (
	"testing"

	"github.com/guestin/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestFinderImpl_FindConfigs(t *testing.T) {
	rootLogger, _ := log.EasyInitConsoleLogger(zap.DebugLevel, zap.DPanicLevel)
	logger := log.NewTaggedZapLogger(rootLogger, "test")
	finder := newConfigFinder(logger)
	_, _ = finder.FindConfigs("./test/config", viper.SupportedExts...)
}
