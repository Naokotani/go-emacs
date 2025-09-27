package main

import (
	"io"
	"os"
)

func fileExists(f string) bool {
	if _, err := os.Stat(f); err == nil {
		return true
	}
	return false
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Sync()
}
