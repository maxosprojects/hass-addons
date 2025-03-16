package main

import (
	"backup-to-s3/logging"
	"backup-to-s3/syncer"
)

func main() {
	snc := syncer.New(logging.New())
	snc.Run()
}
