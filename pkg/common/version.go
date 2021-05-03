package common

import "fmt"

// NAME of the App
var NAME = "integrity"

// SUMMARY of the Version
var SUMMARY = fmt.Sprintf("%s-%s", VERSION, BRANCH)

// BRANCH of the Version
var BRANCH = "dev"

// VERSION of Release
var VERSION = "2.0.1"

// AppVersion --
var AppVersion AppVersionInfo

// AppVersionInfo --
type AppVersionInfo struct {
	Name    string
	Version string
	Branch  string
	Summary string
}

func init() {
	AppVersion = AppVersionInfo{
		Name:    NAME,
		Version: VERSION,
		Branch:  BRANCH,
		Summary: SUMMARY,
	}
}
