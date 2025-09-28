package main

import (
	"io"
	"os"
	"path/filepath"
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

func getDirectoryFiles(path string) ([]string, error) {
	var files []string

	if !fileExists(path) {
		return files, nil
	}

	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}
	return files, nil

}
