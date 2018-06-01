Ambient-Go
=====================

Ambient client library for Go language

  * https://ambidata.io

This library is base on Ambient Python client library
  * https://github.com/AmbientDataInc/ambient-python-lib

```go
package main

import (
  "log"
  "github.com/sakurahilljp/ambient-go"
)

func main() {
  client := ambient.NewClient(1234, "writeky")

  dp := ambient.NewDataPoint()
  dp["d1"] = 1.23

  err := client.Send(dp)
  if err != nil {
	log.Fatal(err)
  }
}
```

## Installation


To install docopt in your `$GOPATH`:

```console
$ go get  -u github.com/sakurahilljp/ambient-go
```

```go
import "github.com/sakurahilljp/ambient-go"
```

## API

### Send

### Read

### GetProp