package redis

import (
	goRedis "github.com/go-redis/redis"
)

type Client struct {
	*goRedis.Client
}

var _redisCli *Client = nil

//goland:noinspection ALL
func Cli() *Client {
	return _redisCli
}
