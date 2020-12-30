package main

import (
    "time"
    "strings"
    "path/filepath"
    "flag"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/user"
)

// Error handling function
func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Hash the file
func hash_file_sha256(filePath string) (string, error) {
    var SHA256String string
    file, err := os.Open(filePath)
    if err != nil {
    	return SHA256String, err
    }
    defer file.Close()
    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
    	return SHA256String, err
    }
    hashInBytes := hash.Sum(nil)[:32]
    SHA256String = hex.EncodeToString(hashInBytes)
    return SHA256String, nil
}

// Create VERSION file and add headings
func createVerFile(verFile string) {
    user, err := user.Current()
    check(err)
    timestamp := time.Now().Format(time.RFC3339)
    f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()
    _, err = f.WriteString("# integrity 1.0.0 output generated on " + timestamp + " by " + user.Name + "\n")
    check(err)
    _, err = f.WriteString("# " + strings.Join(os.Args, " ") + "\n")
    check(err)
    _, err = f.WriteString("# Filename: SHA256\n")
    check(err)
    f.Sync()
}

// Add data for hashed file to VERSION file
func appendVerFile(verFile string, fileName string, SHA256String string) {
    f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    check(err)
    defer f.Close()
    _, err = f.WriteString(fileName + ": " + SHA256String + "\n")
    f.Sync()
}

// Main function
func main() {
    dirPtr := flag.String("d", ".", "Target directory")
    validatePtr := flag.Bool("v", false, "Validate existing VERSION file")
    verPtr := flag.String("c", "UNDEFINED", "Courseware Version Indentifier")
    flag.Parse()
    verVal := *verPtr
    validateVal := *validatePtr
    if (verVal == "UNDEFINED") {
        fmt.Println("You MUST define a Courseware Version Identifier (e.g., 'integrity -c SEC123-21-01')")
        os.Exit(1)
    }

    // Create VERISON FILE ENTRIES
    if validateVal == false {
        fmt.Println("[+] Working directory:", *dirPtr)
        _, err := os.Stat(*dirPtr + "/VERSION-" + *verPtr)
        if (err == nil) {
            fmt.Println("[!] VERSION file already exists! Overwriting!")
            err = os.Remove(*dirPtr + "/VERSION-" + *verPtr)
        }
        createVerFile(*dirPtr + "/VERSION-" + *verPtr)

        err = filepath.Walk(*dirPtr,
        func(path string, info os.FileInfo, err error) error {
            pathCheck, _ := os.Stat(path)
            if ! pathCheck.IsDir() {
                fmt.Println("[+] Processing " + path + "...")
                hash, err := hash_file_sha256(path)
                check(err)
                if err == nil {
                    fileName, err := filepath.Rel(*dirPtr, path)
                    check(err)
                    if (fileName != "VERSION-" + *verPtr) {
                        appendVerFile(*dirPtr + "/VERSION-" + *verPtr, fileName, hash)
                    }
                }
            }
            return nil
        })
    } else {
        //Validate existing VERSION file
        var failed = false
        verFile := *dirPtr + "/VERSION-" + *verPtr
        _, err := os.Stat(verFile)
        check(err)
        verBytes, err := ioutil.ReadFile(verFile)
        check(err)
        vfString := string(verBytes)
        err = filepath.Walk(*dirPtr,
        func(path string, info os.FileInfo, err error) error {
            pathCheck, _ := os.Stat(path)
            if ! pathCheck.IsDir() {
                hash, err := hash_file_sha256(path)
                check(err)
                fileName, err := filepath.Rel(*dirPtr, path)
                check(err)
                if (! strings.Contains(vfString, fileName + ": " + hash)) && (fileName != "VERSION-" + *verPtr) {
                    fmt.Println("[!] Validation failed!")
                    fmt.Println("    File: " + fileName)
                    fmt.Println("    Hash: " + hash)
                    failed = true
                }
            }
            return nil
        })        
        if !failed {
            fmt.Println("[+] Validation succeeded!")
        }
    }
}
