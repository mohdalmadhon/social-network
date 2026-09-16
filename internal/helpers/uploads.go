package helpers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const AVATAR_PATH = "uploads/avatars"
const POSTS_PATH = "uploads/posts"
const AVATARS_GROUPS_PATH = "uploads/groups/avatars"

func SaveUploads(file multipart.File, header *multipart.FileHeader, Type string) (string, error) {
	var path string

	if Type == "post" {
		path = POSTS_PATH
	} else if Type == "avatar" {
		path = AVATAR_PATH
	} else if Type == "group/avatar" {
		path = AVATARS_GROUPS_PATH
	} else {
		return "", fmt.Errorf("invalid upload type")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	extension := filepath.Ext(header.Filename)
	filename := uuid.New().String() + extension
	filePath := filepath.Join(path, filename)

	log.Println("Saving upload to:", filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)

	if err != nil {
		return "", err
	}

	if Type == "avatar" {
		return "avatars/" + filename, nil
	}

	return "posts/" + filename, nil
}
