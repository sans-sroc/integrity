package common

const Filename = "sans-integrity.yml"
const GetFirstDirectory = "get_first"

var IgnoreAlways = []string{
	Filename,
	"*\\.gpg$",
}

var IgnoreOnCreate = []string{
	".DS_Store",
	"desktop.ini",
	".System Volume Information",
}
