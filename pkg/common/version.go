package common

// NAME of the App
var NAME = "integrity"

// SUMMARY of the Version
var SUMMARY = "3.0.1"

// BRANCH of the Version
var BRANCH = "main"

// VERSION of Release
var VERSION = "3.0.1"

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
