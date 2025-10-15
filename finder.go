package kboot

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/guestin/log"
	"go.uber.org/zap"
)

func buildRegex(profile string, exts []string) (*regexp.Regexp, error) {
	pattern := ""
	extReg := make([]string, 0)
	for _, ext := range exts {
		extReg = append(extReg, "\\."+ext)
	}
	if profile != "" {
		pattern = fmt.Sprintf("(\\.|-|_)%s(%s)$", profile, strings.Join(extReg, "|"))
	} else {
		pattern = fmt.Sprintf("^[^\\-._\\n]+(%s)$", strings.Join(extReg, "|"))
	}
	return regexp.Compile(pattern)
}

func findConfigFiles(logger log.ZapLog, paths []string, profile string, exts []string) []string {
	result := make([]string, 0)
	reg, err := buildRegex(profile, exts)
	if err != nil {
		logger.Warn("build finder regexp err", zap.Error(err))
		return result
	}
	logger.Info("try find config file",
		zap.Strings("path", paths),
		zap.String("profile", profile),
		zap.Strings("exts", exts))
	for _, p := range paths {
		files, err := os.ReadDir(p)
		if err != nil {
			logger.Warn("read config dir error", zap.Error(err))
			continue
		}
		for _, f := range files {
			fn := f.Name()
			if reg.MatchString(fn) {
				result = append(result, path.Join(p, fn))
			}
		}
	}
	return result
}
