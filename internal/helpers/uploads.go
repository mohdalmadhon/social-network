package helpers

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	if header.Size > 5*1024*1024 {
		return "", errors.New("image must be no larger than 5 MB")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	config, format, err := image.DecodeConfig(file)
	if err != nil || config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 40_000_000 {
		return "", errors.New("invalid or oversized image")
	}
	extensions := map[string]string{"jpeg": ".jpg", "png": ".png", "gif": ".gif"}
	extension, ok := extensions[format]
	if !ok {
		return "", errors.New("only JPEG, PNG and GIF images are allowed")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
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
