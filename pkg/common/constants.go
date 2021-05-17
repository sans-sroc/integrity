package common

const NameFormat = "[0-9]{3}.[0-9]{2}.[0-9][A-Z]"

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
