package utils

import "os"

func WriteFile(path string) error {
	f, err := os.Create(path)
	f.Close()
	return err
}
