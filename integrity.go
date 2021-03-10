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
func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Hash the file
func hashFileSha256(filePath string) (string, error) {
    var sha256String string
    file, err := os.Open(filePath)
    if err != nil {
    	return sha256String, err
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

// Create VERSION file and add headings
func createVerFile(verFile string) {
    user, err := user.Current()
    check(err)
    timestamp := time.Now().Format(time.RFC3339)
    f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    check(err)
    defer f.Close()
    _, err = f.WriteString("# integrity 1.0.0 output generated on " + timestamp + " by " + user.Name + "\n")
    check(err)
    _, err = f.WriteString("# " + strings.Join(os.Args, " ") + "\n")
    check(err)
    _, err = f.WriteString("# Filename: SHA256\n")
    check(err)
    err = f.Sync()
    check(err)
}

// Add data for hashed file to VERSION file
func appendVerFile(verFile string, fileName string, sha256String string) {
    f, err := os.OpenFile(verFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    check(err)
    defer f.Close()
    _, err = f.WriteString(fileName + ": " + sha256String + "\n")
    check(err)
    err = f.Sync()
    check(err)
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
    if !validateVal {
        fmt.Println("[+] Working directory:", *dirPtr)
        _, err := os.Stat(*dirPtr + "/VERSION-" + *verPtr)
        if (err == nil) {
            fmt.Println("[!] VERSION file already exists! Overwriting!")
            err = os.Remove(*dirPtr + "/VERSION-" + *verPtr)
            check(err)
        }
        createVerFile(*dirPtr + "/VERSION-" + *verPtr)

        err = filepath.Walk(*dirPtr,
        func(path string, info os.FileInfo, err error) error {
            pathCheck, err2 := os.Stat(path)
            check(err2)
            if ! pathCheck.IsDir() {
                fmt.Println("[+] Processing " + path + "...")
                hash, err2 := hashFileSha256(path)
                check(err2)
                if err2 == nil {
                    fileName, err2 := filepath.Rel(*dirPtr, path)
                    check(err2)
                    if (fileName != "VERSION-" + *verPtr) {
                        appendVerFile(*dirPtr + "/VERSION-" + *verPtr, fileName, hash)
                    }
                }
            }
            return nil
        })
        check(err)
    } else {
        // Validate existing VERSION file
        var failed = false
        verFile := *dirPtr + "/VERSION-" + *verPtr
        _, err := os.Stat(verFile)
        check(err)
        verBytes, err := ioutil.ReadFile(verFile)
        check(err)
        vfString := string(verBytes)
        err = filepath.Walk(*dirPtr,
        func(path string, info os.FileInfo, err error) error {
            pathCheck, err2 := os.Stat(path)
            check(err2)
            if ! pathCheck.IsDir() {
                hash, err2 := hashFileSha256(path)
                check(err2)
                fileName, err2 := filepath.Rel(*dirPtr, path)
                check(err2)
                if (! strings.Contains(vfString, fileName + ": " + hash)) && (fileName != "VERSION-" + *verPtr) {
                    fmt.Println("[!] Validation failed!")
                    fmt.Println("    File: " + fileName)
                    fmt.Println("    Hash: " + hash)
                    failed = true
                }
            }
            return nil
        })
        check(err) 
        if !failed {
            fmt.Println("[+] Validation succeeded!")
        }
    }
}
