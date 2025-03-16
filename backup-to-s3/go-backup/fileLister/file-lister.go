package fileLister

import (
	"backup-to-s3/constants"
	"backup-to-s3/options"
	"os"
)

type FileLister interface {
	ListCurrFiles() ([]string, error)
	// GetNewFormat returns the HA backup file name in s3 format, inferred from the new HA backup file name format,
	// in case the provided backup file name appears to be in the new HA backup file name format, and "true" as an
	// indicator the provided filename is indeed in the new HA backup file name format.
	// Otherwise, returns the provided filename unchanged, and "false".
	GetNewFormat(filename string) (string, bool)
}

type fileLister struct {
	opts *options.Options
}

func New(opts *options.Options) FileLister {
	return &fileLister{
		opts: opts,
	}
}

// ListCurrFiles returns the names of the backup files as they appear in HA /backup folder, or as they appear
// in s3 (inferred from the new HA filename format), and values are file names as they appear in HA /backup folder.
func (fl *fileLister) ListCurrFiles() ([]string, error) {
	entries, err := os.ReadDir(fl.opts.BackupDir)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		res = append(res, entry.Name())
	}

	return res, err
}

func (fl *fileLister) GetNewFormat(filename string) (string, bool) {
	matches := constants.NewFilenameFormatRegex.FindStringSubmatch(filename)
	if len(matches) == 5 {
		return matches[1] + ".tar", true
	}
	return filename, false
}
