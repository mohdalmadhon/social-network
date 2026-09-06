package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const AVATAR_PATH = "uploads/avatars"
const POSTS_PATH = "uploads/posts"
const COMMENTS_PATH = "uploads/comments"

func SaveUploads(file multipart.File, header *multipart.FileHeader, Type string) (string, error) {
	var path string

	if Type == "avatar" {
		path = AVATAR_PATH
	} else if Type == "post" {
		path = POSTS_PATH
	} else if Type == "comment" {
		path = COMMENTS_PATH
	} else {
		return "", os.ErrInvalid
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
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
		os.Remove(filePath)
		return "", err
	}

	if Type == "avatar" {
		return "avatars/" + filename, nil
	}
	if Type == "comment" {
		return "comments/" + filename, nil
	}

	return "posts/" + filename, nil
}

func DeleteAvatar(avatarPath string) error {
	if avatarPath == "" {
		return nil
	}

	if filepath.Base(avatarPath) == "default.png" {
		return nil
	}

	filePath := filepath.Join("uploads", filepath.FromSlash(avatarPath))

	return os.Remove(filePath)
}
