package common

const Filename = "sans-integrity.yml"
const FilenameSigned = "sans-integrity.yml.gpg"
const GetFirstDirectory = "get_first"

var IgnoreAlways = []string{
	Filename,
	FilenameSigned,
}

var IgnoreOnCreate = []string{
	".DS_Store",
	"desktop.ini",
	".System Volume Information",
}
