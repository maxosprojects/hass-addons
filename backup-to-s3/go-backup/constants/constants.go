package constants

import "regexp"

var NewFilenameFormatRegex = regexp.MustCompile("^(.+)_([0-9]{4}-[0-9]{2}-[0-9]{2})_([0-9]{2}\\.[0-9]{2})_([0-9]{8})\\.tar$")
