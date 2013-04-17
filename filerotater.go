/*
  PUBLIC DOMAIN STATEMENT
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Source Code file.
  This work is published from the United Kingdom.
*/

// A writer that rotates an underlying file, e.g. for log files
package filerotate

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type FileRotater struct {
	filename     string
	format       string
	rotationtime time.Duration
	guard        sync.Mutex
	file         *os.File
	quit         chan bool
}

// Create a new FileRotater
// filename is the base name of the file.
// format is a date format string that will be applied to the rotation time and appended to
// filename as an extension. Use an empty string to append unixtime seconds.
// rotationtime is the interval between rotations, starting at midnight of current day
func NewFileRotater(filename string, format string, rotationtime time.Duration) (*FileRotater, error) {
	fr := &FileRotater{filename: filename, format: format, rotationtime: rotationtime}
	err := fr.init()
	if err != nil {
		return nil, err
	}
	return fr, nil
}

func (fr *FileRotater) init() error {
	var err error

	quit := make(chan bool)
	n := fr.generateName()

	t := time.Now()
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	interval := time.Since(midnight) % fr.rotationtime

	fr.guard.Lock()
	defer fr.guard.Unlock()
	fr.file, err = os.Create(n)
	if err != nil {
		return err
	}

	go func() {
		timer := time.After(interval)
		for {
			select {

			case <-timer:

				n := fr.generateName()

				fr.guard.Lock()
				err := fr.file.Close()
				if err != nil {
					fr.guard.Unlock()
					return
				}
				fr.file, err = os.Create(n)
				fr.guard.Unlock()
				if err != nil {
					return
				}

				timer = time.After(fr.rotationtime)

			case <-quit:
				return

			}
		}
	}()

	return nil
}

func (fr *FileRotater) generateName() string {
	t := time.Now()
	if fr.format == "" {
		return fmt.Sprintf("%s.%d", fr.filename, t.Unix())
	} else {
		return fmt.Sprintf("%s.%s", fr.filename, t.Format(fr.format))
	}
}

func (fr *FileRotater) Write(p []byte) (n int, err error) {
	fr.guard.Lock()
	defer fr.guard.Unlock()
	return fr.file.Write(p)
}

func (fr *FileRotater) Close() error {
	fr.quit <- true
	fr.guard.Lock()
	defer fr.guard.Unlock()
	return fr.file.Close()

}
