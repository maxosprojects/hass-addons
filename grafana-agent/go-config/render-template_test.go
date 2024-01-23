package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io"
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	optionsPath = "test-data/options.json"
	templatePath = "test-data/template.tmpl"
	renderToPath = "test-data/result.yaml"

	run()

	expectedText := readText("test-data/expected-result.yaml")
	actualText := readText(renderToPath)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expectedText, actualText, true)

	fmt.Println(dmp.DiffPrettyText(diffs))
}

func readText(filename string) string {
	file, err := os.Open(filename)
	noErr(err)
	defer file.Close()

	text, err := io.ReadAll(file)
	noErr(err)
	return string(text)
}
