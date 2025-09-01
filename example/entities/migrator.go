package entities

import "github.com/guestin/kboot/db"

func init() {
	db.SetupMigrateBuilder(func() error {

		return nil
	})
}
