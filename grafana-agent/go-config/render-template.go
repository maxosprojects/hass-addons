package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"text/template"
)

var optionsPath = "/data/options.json"
var templatePath = "/grafana-agent-config.tmpl"
var renderToPath = "/grafana-agent-config.yaml"

func main() {
	run()
}

func run() {
	jsonFile, err := os.Open(optionsPath)
	noErr(err)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var data map[string]any
	err = json.Unmarshal(byteValue, &data)
	noErr(err)

	tmpl, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	noErr(err)

	file, err := os.Create(renderToPath)
	noErr(err)

	err = tmpl.Execute(file, data)
	noErr(err)
}

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
