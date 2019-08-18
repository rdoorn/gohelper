package filehelper

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

func IsDir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func IsWritable(filename string) error {
	if _, err := os.Stat(filename); err == nil {

		// path/to/whatever exists
		file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
		defer file.Close()
		if err != nil {
			if os.IsPermission(err) {
				return err
			}

		}

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		dir := path.Dir(filename)
		if unix.Access(dir, unix.W_OK) != nil {
			return fmt.Errorf("cannot write to lof output %s: directory not writable", dir)
		}
	}
	return nil
}
