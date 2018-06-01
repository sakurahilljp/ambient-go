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
