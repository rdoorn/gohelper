package arg

import (
	"os"
	"path/filepath"
)

func GetName() string {
	return filepath.Base(os.Args[0])
}
