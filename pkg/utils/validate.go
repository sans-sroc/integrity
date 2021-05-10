package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sirupsen/logrus"
)

func ValidateFiles(directory string, version string, parts bool, first bool, jsonOut bool, jsonPretty bool) bool {
	// Be sure to normalize \ to /
	directory = filepath.ToSlash(directory)

	var failed = false
	verFile := ""

	files := []*common.File{}

	versionFileName := fmt.Sprintf("VERSION-%s.txt", version)
	versionPartFileName := fmt.Sprintf("VERSION-%s-part.txt", version)
	versionFirstFileName := fmt.Sprintf("VERSION-%s-first.txt", version)

	if parts {
		verFile = filepath.Join(directory, versionPartFileName)
	} else if first {
		verFile = filepath.Join(directory, versionFirstFileName)
	} else {
		verFile = filepath.Join(directory, versionFileName)
	}
	_, err := os.Stat(verFile)
	check(err, "Cannot find VERSION file")

	verBytes, err := ioutil.ReadFile(verFile)
	check(err, "Cannot read VERSION file")

	vfString := string(verBytes)
	cfileNames := make(map[string]string)
	ofileNames := make(map[string]string)

	re := regexp.MustCompile("(?m)(^[^#][^:]+):(.*)$")
	allFiles := re.FindAllStringSubmatch(vfString, -1)
	for i := 0; i < len(allFiles); i++ {
		ofileNames[filepath.ToSlash(allFiles[i][1])] = strings.TrimSpace(allFiles[i][2])
	}

	currentFiles, err := GetFiles(directory)
	if err != nil {
		logrus.WithError(err).Error("unable to get files")
		return false
	}

	for _, file := range currentFiles {
		if !jsonOut {
			fmt.Println("[+] Validating " + file.Name + " ...")
		}

		hash, err2 := HashFileSha256(file.Path)
		check(err2, "Cannot hash file")

		fileName, err2 := filepath.Rel(directory, file.Path)
		check(err2, "Cannot find file")

		match, _ := regexp.MatchString("VERSION-"+version+".*\\.txt", fileName)
		if !match {
			cfileNames[fileName] = hash
			if _, ok := ofileNames[fileName]; !ok {
				files = append(files, &common.File{
					Name:   fileName,
					Hash:   hash,
					Status: "new",
				})

				if !jsonOut {
					fmt.Println("[!] Validation failed! File has been added!")
					fmt.Println("    File: " + fileName)
					fmt.Println("    Hash: " + hash)
					failed = true
				}
			}

			if _, ok := ofileNames[fileName]; ok {
				if ofileNames[fileName] != hash {
					files = append(files, &common.File{
						Name:   fileName,
						Hash:   hash,
						Status: "failed",
					})

					if !jsonOut {
						fmt.Println("[!] Validation failed! File contents have been modified!")
						fmt.Println("    File: " + fileName)
						fmt.Println("    Hash: " + hash)
						failed = true
					}
				}
			}
		}
	}

	for name, hash := range ofileNames {
		if _, ok := cfileNames[name]; !ok {
			files = append(files, &common.File{
				Name:   name,
				Hash:   hash,
				Status: "missing",
			})

			if !jsonOut {
				fmt.Println("[!] Validation failed! Original file missing!")
				fmt.Println("    File: " + name)
				fmt.Println("    Hash: " + hash)
				failed = true
			}
		}
	}

	if jsonOut {
		var out []byte
		if jsonPretty {
			out, err = json.MarshalIndent(common.FileOutput{Files: files}, "", "    ")
			check(err, "unable to render json")
		} else {
			out, err = json.Marshal(common.FileOutput{Files: files})
			check(err, "unable to render json")
		}
		fmt.Println(string(out))
	} else {
		check(err, "Validation failed")
		fmt.Println("[+] Validation process complete!")
	}

	return failed
}
