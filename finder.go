package kboot

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/guestin/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type configFile struct {
	FileName string
	FilePath string
	FileType string
	Profile  string
}

type configItem struct {
	Name     string
	Default  *configFile
	Profiles []*configFile
}

type configFinder interface {
	FindConfigs(dir string, exts ...string) (map[string]*configItem, error)
}

func newConfigFinder(logger log.ZapLog) configFinder {
	return &configFinderImpl{logger: logger}
}

type configFinderImpl struct {
	logger log.ZapLog
}

func (this *configFinderImpl) FindConfigs(dir string, exts ...string) (map[string]*configItem, error) {
	this.logger.Info("try find config file from",
		zap.String("dir", dir),
		zap.Strings("exts", exts))
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "read config form path %s error", dir)
	}
	pattern := fmt.Sprintf("^([a-zA-Z]+[a-zA-Z0-9]*)[_.-]?([a-zA-Z0-9]*)\\.(%s)$", strings.Join(exts, "|"))
	reg := regexp.MustCompile(pattern)
	result := make(map[string]*configItem)
	for _, f := range files {
		fn := f.Name()
		match := reg.FindStringSubmatch(fn)
		if len(match) == 4 {
			configName := match[1]
			cfg, ok := result[configName]
			if !ok {
				result[configName] = &configItem{
					Name:     configName,
					Profiles: make([]*configFile, 0),
				}
			}
			cfg = result[configName]
			profileFile := &configFile{
				FileName: fn,
				FilePath: path.Join(dir, fn),
				Profile:  match[2],
				FileType: match[3],
			}
			if profileFile.Profile != "" {
				cfg.Profiles = append(cfg.Profiles, profileFile)
			} else {
				cfg.Default = profileFile
			}
		}
	}
	return result, nil
}
