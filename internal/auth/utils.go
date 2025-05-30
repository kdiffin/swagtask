package auth

import (
	"fmt"
	"log"
	"net/textproto"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func validateImage(fileName string, fileSize int64, fileHeader textproto.MIMEHeader) error {
	log.Printf("Uploaded file: %s, Size: %d bytes, MIME Type: %s",
		fileName, fileSize, fileHeader.Get("Content-Type"))

	// Check file type (MIME type)
	contentType := fileHeader.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"} // Define allowed image types

	// Simple check for allowed types
	isAllowedType := false
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			isAllowedType = true
			break
		}
	}
	if !isAllowedType {
		return fmt.Errorf("Unsupported file type: %s. Only %v allowed.", contentType, allowedTypes)
	}

	return nil
}
