package helpers

import (
	"os"
	"path/filepath"
)

func BasePath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
