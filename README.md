# K-BOOT

## Usage

> see example

```bash
go get -u github.com/guestin/kboot
```

```go
package main

import (
	"context"

	"github.com/guestin/kboot"
)

func main() {
	kboot.Bootstrap(context.Background())
}

```
