package apis

import "github.com/guestin/kboot/db"

func Echo(req map[string]interface{}) (map[string]interface{}, error) {
	db.ORM()
	return req, nil
}
