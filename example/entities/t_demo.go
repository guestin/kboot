package entities

import "github.com/guestin/kboot/db"

type Demo struct {
	db.CreatedAt
	db.DeletedAt
	Name string
}

func (this *Demo) TableName() string {
	return "t_demo"
}
