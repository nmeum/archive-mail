package main

import (
	"path/filepath"
)

func isMaildir(name string) bool {
	return name == "new" || name == "cur" || name == "tmp"
}

func getDir(path string) string {
	dir := filepath.Base(filepath.Dir(path))
	if !isMaildir(dir) {
		panic("unexpected non-maildir folder")
	}

	return dir
}
