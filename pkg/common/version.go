package common

// NAME of the App
var NAME = "integrity"

// SUMMARY of the Version, this is using git describe
// Note: This generally gets set to the major version of the app,
// it gets set to the real version during the build process.
var SUMMARY = "3.0.0"

// BRANCH of the Version
var BRANCH = "main"

// VERSION of Release
// Note: This generally gets set to the major version of the app,
// it gets set to the real version during the build process.
var VERSION = "3.0.0"

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
