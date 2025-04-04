package main

import (
	"encoding/json"
	"os"
)

func LogStart(d StartJSON) {
	b, err := json.Marshal(&d)
	if err != nil {
		LogError("json Marshal" + err.Error())
		return
	}

	// Write data
	if _, err = os.Stdout.Write(append(b, '\n')); err != nil {
		LogError("write data" + err.Error())
		return
	}
}

func LogError(d string) {
	// Write data
	_, _ = os.Stderr.WriteString(d + "\n")
}
