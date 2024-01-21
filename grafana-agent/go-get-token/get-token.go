package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var optionsPath = "/data/options.json"

func main() {
	run()
}

func run() {
	jsonFile, err := os.Open(optionsPath)
	noErr(err)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	noErr(err)

	arg := os.Args[1]
	fmt.Println(result[arg])
}

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
