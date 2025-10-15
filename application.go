package kboot

import "time"

type Application interface {
	GetAppName() string
	GetTimezone() *time.Location
}
