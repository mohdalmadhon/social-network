package validation

import (
	"mime/multipart"
	"path/filepath"
	"social/internal/models"
	"strconv"
	"strings"
)

type ValidationResult struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidatePost(data *models.RegsiterPost, image *multipart.FileHeader) ValidationResult {
	if data.Content == "" {
		return ValidationResult{
			Field:   "content",
			Message: "content cannot be empty",
		}
	}

	if len([]rune(data.Content)) > 1000 {
		return ValidationResult{
			Field:   "content",
			Message: "content cannot be more than 1000 characters",
		}
	}

	if data.AllowComments != 0 && data.AllowComments != 1 {
		return ValidationResult{
			Field:   "allowComments",
			Message: "invalid allow comments code",
		}
	}

	if data.GroupID < -1 {
		return ValidationResult{
			Field:   "groupID",
			Message: "invalid group code",
		}
	}

	if data.Location != "" {
		locationParts := strings.Split(data.Location, ":")

		if len(locationParts) != 3 ||
			strings.TrimSpace(locationParts[0]) == "" ||
			strings.TrimSpace(locationParts[1]) == "" ||
			strings.TrimSpace(locationParts[2]) == "" {
			return ValidationResult{
				Field:   "location",
				Message: "invalid location",
			}
		}

		if _, err := strconv.ParseFloat(strings.TrimSpace(locationParts[1]), 64); err != nil {
			return ValidationResult{
				Field:   "location",
				Message: "invalid location",
			}
		}

		if _, err := strconv.ParseFloat(strings.TrimSpace(locationParts[2]), 64); err != nil {
			return ValidationResult{
				Field:   "location",
				Message: "invalid location",
			}
		}
	}

	for _, id := range data.PeopleTagged {
		if id <= 0 {
			return ValidationResult{
				Field:   "tags",
				Message: "invalid tagged people",
			}
		}
	}

	if image != nil {
		allowedExtensions := map[string]bool{
			".png":  true,
			".jpg":  true,
			".jpeg": true,
			".gif":  true,
		}

		extension := strings.ToLower(filepath.Ext(image.Filename))

		if !allowedExtensions[extension] {
			return ValidationResult{
				Field:   "image",
				Message: "image must be png, jpg or gif",
			}
		}

		if image.Size <= 0 {
			return ValidationResult{
				Field:   "image",
				Message: "invalid image",
			}
		}
	}

	return ValidationResult{}
}
