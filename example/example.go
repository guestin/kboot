package main

import (
	"context"

	"github.com/guestin/kboot"
	"github.com/guestin/kboot/db"
)
import _ "github.com/guestin/kboot/example/apis"

func main() {
	kboot.Bootstrap(context.Background())
	db.ORM("")
}
