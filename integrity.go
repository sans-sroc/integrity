package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var rel_ver = "1.0.7"

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

// Normalize slashes
func normalizeSlashes(filePath string) string {
	filePath = strings.Replace(filePath, "\\", "/", -1)
	return filePath
}

// Create VERSION file and add headings
func createVerFile(verFile string) {
	user, err := user.Current()
	check(err, "Cannot determine user")
	timestamp := time.Now().Format(time.RFC3339)
	f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot open file")
	defer f.Close()
	_, err = f.WriteString("# integrity " + rel_ver + " output generated on " + timestamp + " by " + user.Name + "\n")
	check(err, "Cannot write to file")
	_, err = f.WriteString("# " + strings.Join(os.Args, " ") + "\n")
	check(err, "Cannot write to file")
	_, err = f.WriteString("# Filename: SHA256\n")
	check(err, "Cannot write to file")
	err = f.Sync()
	check(err, "Cannot write to file")
}

// Add data for hashed file to VERSION file
func appendVerFile(verFile string, fileName string, sha256String string, dirVal string) {
	// Main VERSION file
	f, err := os.OpenFile(verFile+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, "Cannot open file")
	defer f.Close()
	_, err = f.WriteString(fileName + ": " + sha256String + "\n")
	check(err, "Cannot write to file")
	err = f.Sync()
	check(err, "Cannot write to file")

	// Part VERSION file
	dirVal = normalizeSlashes(dirVal)
	_, err1 := os.Stat(dirVal + "/get_first")
	if err1 == nil {
		match, _ := regexp.MatchString("get[-_]first", fileName)
		if !match {
			f, err := os.OpenFile(verFile+"-part.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err, "Cannot open file")
			defer f.Close()
			fileName = normalizeSlashes(fileName)
			_, err = f.WriteString(fileName + ": " + sha256String + "\n")
			check(err, "Cannot write to file")
			err = f.Sync()
			check(err, "Cannot write to file")
		}

		// First VERSION file
		if match {
			f, err := os.OpenFile(verFile+"-first.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err, "Cannot open file")
			defer f.Close()
			fileName = normalizeSlashes(fileName)
			_, err = f.WriteString(filepath.Base(fileName) + ": " + sha256String + "\n")
			check(err, "Cannot write to file")
			err = f.Sync()
			check(err, "Cannot write to file")
		}
	}

}

func validateFiles(directory string, version string, parts bool, first bool, json bool) bool {
	var failed = false
	jsonFirst := true
	jsonOutput := "{\n\t\"files\": [\n"
	verFile := ""
	if parts {
		verFile = directory + "/VERSION-" + version + "-part.txt"
	} else if first {
		verFile = directory + "/VERSION-" + version + "-first.txt"
	} else {
		verFile = directory + "/VERSION-" + version + ".txt"
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
		ofileNames[allFiles[i][1]] = strings.TrimSpace(allFiles[i][2])
	}
	err = filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			pathCheck, err2 := os.Stat(path)
			check(err2, "Cannot find file")
			if !pathCheck.IsDir() {
				hash, err2 := hashFileSha256(path)
				check(err2, "Cannot hash file")
				fileName, err2 := filepath.Rel(directory, path)
				check(err2, "Cannot find file")
				match, _ := regexp.MatchString("VERSION-"+version+".*\\.txt", fileName)
				if !match {
					cfileNames[fileName] = hash
					if _, ok := ofileNames[fileName]; !ok {
						if !json {
							fmt.Println("[!] Validation failed! File has been added!")
							fmt.Println("    File: " + fileName)
							fmt.Println("    Hash: " + hash)
							failed = true
						} else {
							if !jsonFirst {
								jsonOutput += ",\n"
							}
							jsonOutput = jsonOutput + "\t\t{\n\t\t\t\"fileName\": \"" + fileName + "\",\n\t\t\t\"hash\": \"" + hash + "\",\n\t\t\t\"status\": \"new\"\n\t\t}"
							jsonFirst = false
						}
					}
					if _, ok := ofileNames[fileName]; ok {
						if ofileNames[fileName] != hash {
							if !json {
								fmt.Println("[!] Validation failed! File contents have been modified!")
								fmt.Println("    File: " + fileName)
								fmt.Println("    Hash: " + hash)
								failed = true
							} else {
								if !jsonFirst {
									jsonOutput += ",\n"
								}
								jsonOutput = jsonOutput + "\t\t{\n\t\t\t\"fileName\": \"" + fileName + "\",\n\t\t\t\"hash\": \"" + hash + "\",\n\t\t\t\"status\": \"failed\"\n\t\t}"
								jsonFirst = false
							}
						}
					}
				}
			}
			return nil
		})
	for name, hash := range ofileNames {
		if _, ok := cfileNames[name]; !ok {
			if !json {
				fmt.Println("[!] Validation failed! Original file missing!")
				fmt.Println("    File: " + name)
				fmt.Println("    Hash: " + hash)
				failed = true
			} else {
				if !jsonFirst {
					jsonOutput += ",\n"
				}
				jsonOutput = jsonOutput + "\t\t{\n\t\t\t\"fileName\": \"" + name + "\",\n\t\t\t\"hash\": \"" + hash + "\",\n\t\t\t\"status\": \"missing\"\n\t\t}"
				jsonFirst = false
			}
		}
	}
	if json {
		jsonOutput += "\n\t]\n}"
		fmt.Println(jsonOutput)
	} else {
		check(err, "Validation failed")
		fmt.Println("[+] Validation process complete!")
	}
	return failed
}

// Main function
func main() {
	dirPtr := flag.String("d", ".", "Target directory")
	verPtr := flag.String("c", "UNDEFINED", "Courseware Version Indentifier")
	jsonPtr := flag.Bool("j", false, "Output data in JSON format to stdout")
	validatePtr := flag.Bool("v", false, "Validate existing VERSION file")
	partsPtr := flag.Bool("p", false, "Validate the VERSION-part.txt file")
	firstPtr := flag.Bool("f", false, "Validate the VERSION-first.txt file")
	appVerPtr := flag.Bool("version", false, "Print the integrity version and exit")
	flag.Parse()
	jsonVal := *jsonPtr
	jsonOutput := ""
	jsonFirst := true
	verVal := *verPtr
	validateVal := *validatePtr
	partsVal := *partsPtr
	firstVal := *firstPtr
	appVer := *appVerPtr
	dirVal := *dirPtr
	if appVer {
		fmt.Printf("integrity version %s \n", rel_ver)
		os.Exit(0)
	}
	if verVal == "UNDEFINED" {
		fmt.Print("Enter a Courseware Version Identifier (e.g., 'SEC123-21-01'): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		check(err, "Invalid input")
		input = strings.Replace(input, "\n", "", -1)
		*verPtr = input
	}

	// Process Files
	if !validateVal && !partsVal && !firstVal {
		if !jsonVal {
			fmt.Println("[+] Working directory:", *dirPtr)
			_, err := os.Stat(*dirPtr + "/VERSION-" + *verPtr + ".txt")
			if err == nil {
				fmt.Println("[!] VERSION file already exists! Overwriting!")
				err = os.Remove(*dirPtr + "/VERSION-" + *verPtr + ".txt")
				check(err, "Cannot delete VERSION file")
			}
			_, err1 := os.Stat(*dirPtr + "/VERSION-" + *verPtr + "-part.txt")
			if err1 == nil {
				err = os.Remove(*dirPtr + "/VERSION-" + *verPtr + "-part.txt")
				check(err, "Cannot delete VERSION file")
			}
			_, err2 := os.Stat(*dirPtr + "/VERSION-" + *verPtr + "-first.txt")
			if err2 == nil {
				err = os.Remove(*dirPtr + "/VERSION-" + *verPtr + "-first.txt")
				check(err, "Cannot delete VERSION file")
			}
			_, err3 := os.Stat(*dirPtr + "/get_first")
			if err3 == nil {
				createVerFile(*dirPtr + "/VERSION-" + *verPtr + "-part.txt")
				createVerFile(*dirPtr + "/VERSION-" + *verPtr + "-first.txt")
			}
			createVerFile(*dirPtr + "/VERSION-" + *verPtr + ".txt")
		} else {
			jsonOutput = "{\n\t\"files\": [\n"
		}
		err := filepath.Walk(*dirPtr,
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
					if (!strings.Contains(fileName, "VERSION-"+*verPtr)) && (!jsonVal) {
						appendVerFile(*dirPtr+"/VERSION-"+*verPtr, fileName, hash, dirVal)
					}
					if (!strings.Contains(fileName, "VERSION-"+*verPtr)) && (jsonVal) {
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
		// Validate existing VERSION file(s)
		failed := validateFiles(*dirPtr, *verPtr, *partsPtr, *firstPtr, jsonVal)
		if !jsonVal {
			if failed {
				fmt.Println("[!] Result: FAIL!")
			} else {
				fmt.Println("[+] Result: SUCCESS!")
			}
		}
	}
}
