package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sirupsen/logrus"
)

// Create VERSION file and add headings
func CreateVerFile(verFile, username string) error {
	timestamp := time.Now().Format(time.RFC3339)
	f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.WithError(err).Error("cannot open file")
		return err
	}
	defer f.Close()

	_, err = f.WriteString("# integrity " + common.AppVersion.Version + " output generated on " + timestamp + " by " + username + "\n")
	if err != nil {
		logrus.WithError(err).Error("cannot write to file")
		return err
	}

	_, err = f.WriteString("# " + strings.Join(os.Args, " ") + "\n")
	if err != nil {
		logrus.WithError(err).Error("cannot write to file")
		return err
	}

	_, err = f.WriteString("# Filename: SHA256\n")
	if err != nil {
		logrus.WithError(err).Error("cannot write to file")
		return err
	}

	err = f.Sync()
	if err != nil {
		logrus.WithError(err).Error("cannot write to file")
		return err
	}

	return nil
}

// Add data for hashed file to VERSION file
func AppendVerFile(verFile, verPartFile, verFirstFile, fileName, sha256String, dirVal string, getFirstExists, getFirstIsEmpty bool) {
	// Main VERSION file
	f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot open file")
	defer f.Close()

	fileName = NormalizeSlashes(fileName)

	_, err = f.WriteString(fileName + ": " + sha256String + "\n")
	check(err, "Cannot write to file")

	if getFirstExists && !getFirstIsEmpty {
		if strings.Contains(fileName, "get_first") {
			f, err := os.OpenFile(verFirstFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err, "Cannot open file")
			defer f.Close()

			fileName = NormalizeSlashes(fileName)

			_, err = f.WriteString(filepath.Base(fileName) + ": " + sha256String + "\n")
			check(err, "Cannot write to file")
		} else {
			f, err := os.OpenFile(verPartFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err, "Cannot open file")
			defer f.Close()

			fileName = NormalizeSlashes(fileName)

			_, err = f.WriteString(fileName + ": " + sha256String + "\n")
			check(err, "Cannot write to file")
		}
	}
}

// Borrowed from https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func IsDirectoryEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
