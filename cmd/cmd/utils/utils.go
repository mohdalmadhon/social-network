package utils

import (
	"io"
	"os"
)

func InitUploadFolders() error {
	folders := []string{
		"uploads",
		"uploads/avatars",
		"uploads/posts",
	}

	for _, folder := range folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}

func CopyDefaultAvatar() error {
	source := "images/avatar/default.png"
	destination := "uploads/avatars/default.png"

	if _, err := os.Stat(destination); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll("uploads/avatars", 0755); err != nil {
		return err
	}

	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
