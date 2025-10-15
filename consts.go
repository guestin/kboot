package kboot

import "time"

var DefaultAppTz = time.FixedZone("Asia/Shanghai", 0)

const (
	LoggerTag = "kboot"

	DefaultLogLevel        = "debug"
	DefaultConfigName      = "application"
	DefaultConfigFilePath  = "./config"
	DefaultConfigEnvPrefix = ""

	CfgKeyProfilesActive = "kboot.profiles.active"
	CfgKeyAppTz          = "app.timezone"
	CfgKeyAppLogLevel    = "app.log.level"
)
