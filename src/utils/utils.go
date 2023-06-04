package utils

import (
	"os"
	"time"
)

func OpenOrCreate(path string, flag int) (file *os.File, err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
	file, err = os.OpenFile(path, flag, 0644)
	return file, err
}

func IsoDateTime(time time.Time) string {
	return time.Format("2006-01-02T15:04:05Z07:00")
}
