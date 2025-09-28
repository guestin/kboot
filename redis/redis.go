package redis

import (
	goRedis "github.com/redis/go-redis/v9"
)

type Client struct {
	goRedis.UniversalClient
}

var _redisCli *Client = nil

//goland:noinspection ALL
func Cli() *Client {
	return _redisCli
}
