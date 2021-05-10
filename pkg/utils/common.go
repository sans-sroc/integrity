package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sirupsen/logrus"
)

// Error handling function
func check(e error, m string) {
	if e != nil {
		logrus.WithError(e).Fatal(m)
	}
}

// Error handling function
func Check(e error, m string) {
	check(e, m)
}

func GetFiles(dir string) (files []*common.File, err error) {
	if err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			pathCheck, err := os.Stat(path)
			if err != nil {
				logrus.WithError(err).Error("Cannot process file")
				return err
			}

			if !pathCheck.IsDir() {
				fileName, err := filepath.Rel(dir, path)
				if err != nil {
					logrus.WithError(err).Error("Cannot determine file path")
					return err
				}

				if strings.HasPrefix(fileName, "VERSION-") && strings.HasSuffix(fileName, ".txt") {
					logrus.WithField("filename", fileName).Debug("omitted file")
					return nil
				}

				logrus.WithField("filename", fileName).Debug("found file")

				// Both name and path must be ToSlash because the Name is what
				// is ultimately written to the versioning file
				files = append(files, &common.File{
					Name: filepath.ToSlash(fileName),
					Path: filepath.ToSlash(path),
				})
			}

			return nil
		},
	); err != nil {
		logrus.WithError(err).Error("unable to obtain a file list")
		return nil, err
	}

	return files, nil
}
