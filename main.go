package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/a15y87/udslogbeat/beater"
)

func main() {
	err := beat.Run("udslogbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
