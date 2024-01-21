package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudkucooland/go-kasa"
	"log"
	"os"
	"strings"
	"time"
)

type ConfigEntries struct {
	Data struct {
		Entries []*ConfigEntry `json:"entries"`
	} `json:"data"`
}

type ConfigEntry struct {
	EntryId string `json:"entry_id"`
	Domain  string `json:"domain"`
	Title   string `json:"title"`
	Data    struct {
		Host string `json:"host"`
	} `json:"data"`
	UniqueId   string  `json:"unique_id"`
	DisabledBy *string `json:"disabled_by"`
}

type Result struct {
	config *ConfigEntry
	err    error
}

var baseDir = "/homeassistant/.storage/"
var configsFile = "core.config_entries"

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Addon ready")
	for scanner.Scan() {
		err := runCommand(scanner.Text())
		if err != nil {
			logError(err)
		}
	}

	if err := scanner.Err(); err != nil {
		logError(err)
	}
}

func runCommand(cmdString string) error {
	var cmd map[string]string

	err := json.Unmarshal([]byte(cmdString), &cmd)
	if err != nil {
		return err
	}

	if cmd["cmd"] != "reboot" {
		return errors.New(fmt.Sprintf("Unknown command: '%v', full cmd string: '%s'", cmd["cmd"], cmdString))
	}

	configs, err := getData(getFilename(configsFile), &ConfigEntries{})
	if err != nil {
		return err
	}

	tplinkEntries := filterTplinkDevices(configs)

	displayData(tplinkEntries)

	err = rebootAll(tplinkEntries)
	if err != nil {
		return err
	}

	return nil
}

func filterTplinkDevices(configs *ConfigEntries) []*ConfigEntry {
	var tplinkEntries []*ConfigEntry
	for _, entry := range configs.Data.Entries {
		if entry.Domain == "tplink" {
			tplinkEntries = append(tplinkEntries, entry)
		}
	}
	return tplinkEntries
}

func rebootAll(configs []*ConfigEntry) error {
	resultsChan := make(chan *Result, len(configs)*2)

	enabled := 0

	for _, config := range configs {
		if config.DisabledBy != nil {
			logInfo(fmt.Sprintf("Skipping disabled kasa device '%s:%s'", config.Title, config.Data.Host))
			continue
		}

		go reboot(config, resultsChan)
		enabled++
		time.Sleep(1 * time.Second)
	}

	var failed []string

	for i := 0; i < enabled; i++ {
		res := <-resultsChan
		if res.err != nil {
			failed = append(failed, res.config.Title)
		}
	}

	if len(failed) == 0 {
		logInfo(fmt.Sprintf("All kasa devices rebooted successfully on %s", getNow()))
		return nil
	}

	return errors.New(fmt.Sprintf("Kasa failed to reboot some devices on %s: %s",
		getNow(), strings.Join(failed, ", ")))
}

func reboot(config *ConfigEntry, resultsChan chan *Result) {
	if config.Data.Host == "" {
		errString := fmt.Sprintf("Kasa device '%s' has no IP address", config.Title)
		resultsChan <- makeResult(config, errors.New(errString))
		return
	}

	dev, err := kasa.NewDevice(config.Data.Host)
	if err != nil {
		resultsChan <- makeResult(config, err)
		return
	}

	waitChan := make(chan *Result, 1)

	go func() {
		defer close(waitChan)

		_, err2 := dev.GetSettings()
		waitChan <- makeResult(config, err2)
	}()

	select {
	case res := <-waitChan:
		if res.err != nil {
			resultsChan <- res
			return
		}
	case <-time.After(5 * time.Second):
		resultsChan <- makeResult(config, errors.New(fmt.Sprintf("'Kasa device %s' timed out", config.Title)))
		return
	}

	err = dev.Reboot()
	if err != nil {
		resultsChan <- makeResult(config, err)
		return
	}

	resultsChan <- makeResult(config, nil)
}

func makeResult(config *ConfigEntry, err error) *Result {
	return &Result{
		err:    err,
		config: config,
	}
}

func displayData(configs []*ConfigEntry) {
	res := "tplink entries:\n"
	for _, entry := range configs {
		status := "enabled"
		if entry.DisabledBy != nil {
			status = "disabled"
		}
		res += fmt.Sprintf("  %s: %s [%s]\n", entry.Title, status, entry.Data.Host)
	}

	logInfo(res)
}

func getNow() string {
	return time.Now().In(time.Local).Format("2006-01-02 15:04:05")
}

func getFilename(filename string) string {
	return baseDir + filename
}

func getData[T any](filename string, dest *T) (*T, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func logError(err error) {
	log.Printf("[kasa-restart addon] ERROR %v", err)
}

func logInfo(msg string) {
	log.Printf("[kasa-restart addon] INFO %s", msg)
}

func noErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
