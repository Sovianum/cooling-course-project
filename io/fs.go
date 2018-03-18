package io

import "os"

func PrepareDirectories(paths ...string) {
	for _, path := range paths {
		os.MkdirAll(path, os.ModePerm)
	}
}
