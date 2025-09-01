package kboot

import "time"

const LoggerTag = "bootx"

const DefaultAppName = "TODO"

var DefaultAppTz = time.FixedZone("Asia/Shanghai", 0)

const DefaultLogLevel = "debug"
const DefaultConfigName = "application"
const DefaultConfigFileType = "toml"
const DefaultConfigFilePath = "./config"
const DefaultConfigEnvPrefix = "KT"

const CfgKeyAppName = "app.name"
const CfgKeyAppTz = "app.timezone"
const CfgKeyAppLogLevel = "app.log.level"
