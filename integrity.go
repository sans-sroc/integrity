package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// Error handling function
func check(e error, m string) {
	if e != nil {
		fmt.Println(m)
		panic(e)
	}
}

// Hash the file
func hashFileSha256(filePath string) (string, error) {
	var sha256String string
	file, err := os.Open(filePath)
	check(err, "Error opening file")

	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return sha256String, err
	}
	hashInBytes := hash.Sum(nil)[:32]
	sha256String = hex.EncodeToString(hashInBytes)
	return sha256String, nil
}

// Create VERSION file and add headings
func createVerFile(verFile string) {
	user, err := user.Current()
	check(err, "Cannot determine user")
	timestamp := time.Now().Format(time.RFC3339)
	f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot open file")
	defer f.Close()
	_, err = f.WriteString("# integrity 1.0.0 output generated on " + timestamp + " by " + user.Name + "\n")
	check(err, "Cannot write to file")
	_, err = f.WriteString("# " + strings.Join(os.Args, " ") + "\n")
	check(err, "Cannot write to file")
	_, err = f.WriteString("# Filename: SHA256\n")
	check(err, "Cannot write to file")
	err = f.Sync()
	check(err, "Cannot write to file")
}

// Add data for hashed file to VERSION file
func appendVerFile(verFile string, fileName string, sha256String string) {
	f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot open file")
	defer f.Close()
	_, err = f.WriteString(fileName + ": " + sha256String + "\n")
	check(err, "Cannot write to file")
	err = f.Sync()
	check(err, "Cannot write to file")
}

// Main function
func main() {
	dirPtr := flag.String("d", ".", "Target directory")
	validatePtr := flag.Bool("v", false, "Validate existing VERSION file")
	verPtr := flag.String("c", "UNDEFINED", "Courseware Version Indentifier")
	jsonPtr := flag.Bool("j", false, "Output data in JSON format to stdout")
	flag.Parse()
	jsonVal := *jsonPtr
	jsonOutput := ""
	jsonFirst := true
	verVal := *verPtr
	validateVal := *validatePtr
	if verVal == "UNDEFINED" {
		fmt.Println("You MUST define a Courseware Version Identifier (e.g., 'integrity -c SEC123-21-01')")
		os.Exit(1)
	}

	// Process Files
	if !validateVal {
		if !jsonVal {
			fmt.Println("[+] Working directory:", *dirPtr)
		}
		_, err := os.Stat(*dirPtr + "/VERSION-" + *verPtr)
		if err == nil {
			fmt.Println("[!] VERSION file already exists! Overwriting!")
			err = os.Remove(*dirPtr + "/VERSION-" + *verPtr)
			check(err, "Cannot delete VERSION file")
		}

		if !jsonVal {
			createVerFile(*dirPtr + "/VERSION-" + *verPtr)
		} else {
			jsonOutput = "{\n\t\"files\": [\n"
		}
		err = filepath.Walk(*dirPtr,
			func(path string, info os.FileInfo, err error) error {
				pathCheck, err2 := os.Stat(path)
				check(err2, "Cannot process file")
				if !pathCheck.IsDir() {
					if !jsonVal {
						fmt.Println("[+] Processing " + path + "...")
					}
					hash, err2 := hashFileSha256(path)
					check(err2, "Cannot hash file")
					fileName, err2 := filepath.Rel(*dirPtr, path)
					check(err2, "Cannot determine file path")
					if fileName != "VERSION-"+*verPtr && !jsonVal {
						appendVerFile(*dirPtr+"/VERSION-"+*verPtr, fileName, hash)
					}
					if fileName != "VERSION-"+*verPtr && jsonVal {
						if !jsonFirst {
							jsonOutput += ",\n"
						}
						jsonOutput = jsonOutput + "\t\t{\n\t\t\t\"fileName\": \"" + fileName + "\",\n\t\t\t\"hash\": \"" + hash + "\"\n\t\t}"
						jsonFirst = false
					}

				}
				return nil
			})
		check(err, "Validation failed")
		if jsonVal {
			jsonOutput += "\n\t]\n}"
			fmt.Println(jsonOutput)
		}
	} else {
		// Validate existing VERSION file
		var failed = false
		verFile := *dirPtr + "/VERSION-" + *verPtr
		_, err := os.Stat(verFile)
		check(err, "Cannot find VERSION file")
		verBytes, err := ioutil.ReadFile(verFile)
		check(err, "Cannot read VERSION file")
		vfString := string(verBytes)
		err = filepath.Walk(*dirPtr,
			func(path string, info os.FileInfo, err error) error {
				pathCheck, err2 := os.Stat(path)
				check(err2, "Cannot find file")
				if !pathCheck.IsDir() {
					hash, err2 := hashFileSha256(path)
					check(err2, "Cannot hash file")
					fileName, err2 := filepath.Rel(*dirPtr, path)
					check(err2, "Cannot find file")
					if (!strings.Contains(vfString, fileName+": "+hash)) && (fileName != "VERSION-"+*verPtr) {
						fmt.Println("[!] Validation failed!")
						fmt.Println("    File: " + fileName)
						fmt.Println("    Hash: " + hash)
						failed = true
					}
				}
				return nil
			})
		check(err, "Validation failed")
		fmt.Println("[+] Validation process complete!")
	}
}
