package main

import (
	"context"

	"github.com/guestin/kboot"
)
import _ "github.com/guestin/kboot/example/apis"

func main() {
	kboot.Bootstrap(context.Background())
}
