# K-BOOT

## Usage


> see [example](https://github.com/gustin/kboot-examples)

```bash
go get -u github.com/guestin/kboot
```

```go
package main

import (
	"context"

	"github.com/guestin/kboot"
)

type ExampleApplication struct {
}

func (this *ExampleApplication) GetAppName() string {
	return "example"
}

func (this *ExampleApplication) GetTimezone() *time.Location {
	return time.Local
}

func main() {
	kboot.Bootstrap(context.Background(),&ExampleApplication{})
}

```
