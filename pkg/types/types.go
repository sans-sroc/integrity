package types

type File struct {
	Name   string `json:"file" yaml:"file"`
	Path   string `json:"-" yaml:"-"`
	Hash   string `json:"hash" yaml:"hash"`
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}
