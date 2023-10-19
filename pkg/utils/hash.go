package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Hash the file
func HashFileSha256(filePath string) (string, error) {
	var sha256String string
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return sha256String, err
	}

	hashInBytes := hash.Sum(nil)[:32]
	sha256String = hex.EncodeToString(hashInBytes)

	return sha256String, nil
}
