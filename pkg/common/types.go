package common

type File struct {
	Name   string `json:"fileName"`
	Hash   string `json:"hash"`
	Status string `json:"status,omitempty"`
}

type FileOutput struct {
	Files []File `json:"files"`
}
