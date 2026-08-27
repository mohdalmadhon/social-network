package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const AVATAR_PATH = "uploads/avatars"

func SaveUploads(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll(AVATAR_PATH, 0755); err != nil {
		return "", err
	}

	extension := filepath.Ext(header.Filename)
	filename := uuid.New().String() + extension

	filePath := filepath.Join(AVATAR_PATH, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "avatars/" + filename, nil
}