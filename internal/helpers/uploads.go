package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const AVATAR_PATH = "uploads/posts"
const POSTS_PATH = "uploads/avatars"

func SaveUploads(file multipart.File, header *multipart.FileHeader, Type string) (string, error) {
	if err := os.MkdirAll(AVATAR_PATH, 0755); err != nil {
		return "", err
	}

	var path string
	if Type == "post" {
		path = POSTS_PATH
	} else if Type == "avatar" {
		path = AVATAR_PATH
	}
	
	extension := filepath.Ext(header.Filename)
	filename := uuid.New().String() + extension

	filePath := filepath.Join(path, filename)

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
