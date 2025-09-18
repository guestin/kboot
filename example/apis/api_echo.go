package apis

import (
	"github.com/guestin/kboot/db"
	"github.com/guestin/kboot/redis"
)

func Echo(req map[string]interface{}) (map[string]interface{}, error) {
	db.ORM()
	redis.Cli()
	return req, nil
}
