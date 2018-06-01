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
  dp["d1"] = ReadTemperature()
  dp["d2"] = ReadPressure()
  dp["d3"] = ReadHumidity()

  err := client.Send(dp)
  if err != nil {
    log.Fatal(err)
  }
}
```

## Installation


To install ambient-go in your `$GOPATH`:

```console
$ go get  -u github.com/sakurahilljp/ambient-go
```

```go
import "github.com/sakurahilljp/ambient-go"
```

## API

GoDoc: https://godoc.org/github.com/sakurahilljp/ambient-go

### Send

Send a signle point with automatic timestamp

```go
client := ambient.NewClient(1234, "writekey")
	
dp := ambient.NewDataPoint()
dp["d1"] = 19.2
dp["d2"] = 21.3
  
client.Send(dp)
```

Send multiple points with explicit timestamp

```go
client := ambient.NewClient(1234, "writekey")
	
t1 := time.Now()
dp1 := ambient.NewDataPoint(t1)
dp1["d1"] = 1.23

t2 := time.Now()
dp2 := ambient.NewDataPoint(t2)
dp2["d1"] = 2.34

c.Send(dp1, dp2)
```


### Read

Specifies data points with count and skip
```go
values, err := client.Read(ambient.Count(100))
values, err := client.Read(ambient.Count(100), ambient.Skip(100))
```

Read data points at a specified date
```go
values, err := client.Read(ambient.Date(time.Now())
```

Read data points in a specified time ranage.
```go
ent := time.Now()
start := end.Add(-time.Hour * 24)

values, err := client.Read(ambient.Range(start, end))
```

### GetProp

Get properties of a channel.
```go
prop, err := client.GetProp()
```