package utils

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
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
func AppendVerFile(verFile, verPartFile, verFirstFile, fileName, sha256String, dirVal string) {
	// Main VERSION file
	f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot open file")
	defer f.Close()
	fileName = NormalizeSlashes(fileName)
	_, err = f.WriteString(fileName + ": " + sha256String + "\n")
	check(err, "Cannot write to file")

	// Part VERSION file
	dirVal = NormalizeSlashes(dirVal)
	_, err1 := os.Stat(path.Join(dirVal, "get_first"))
	if err1 == nil {
		match, _ := regexp.MatchString("get[-_]first", fileName)
		if !match {
			f, err := os.OpenFile(verPartFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err, "Cannot open file")
			defer f.Close()
			fileName = NormalizeSlashes(fileName)
			_, err = f.WriteString(fileName + ": " + sha256String + "\n")
			check(err, "Cannot write to file")
		}

		// First VERSION file
		if match {
			f, err := os.OpenFile(verFirstFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err, "Cannot open file")
			defer f.Close()
			fileName = NormalizeSlashes(fileName)
			_, err = f.WriteString(filepath.Base(fileName) + ": " + sha256String + "\n")
			check(err, "Cannot write to file")
		}
	}

}
