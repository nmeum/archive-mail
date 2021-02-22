package main

import (
	"errors"
	"os"
	"path/filepath"
)

func isMaildirFn(name string) bool {
	return name == "new" || name == "cur" || name == "tmp"
}

func isValidMaildir(dir string) bool {
	for _, fn := range []string{"new", "cur", "tmp"} {
		_, err := os.Stat(filepath.Join(dir, fn))
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}

	return true
}

func getDir(path string) string {
	dir := filepath.Base(filepath.Dir(path))
	if !isMaildirFn(dir) {
		panic("unexpected non-maildir folder")
	}

	return dir
}
