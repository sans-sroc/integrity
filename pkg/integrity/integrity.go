package integrity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/sans-sroc/integrity/pkg/common"
	"github.com/sans-sroc/integrity/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Metadata struct {
	Name      string    `json:"name" yaml:"name"`
	CreatedBy string    `json:"created_by" yaml:"created_by"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	Version   string    `json:"version" yaml:"version"`
	Directory string    `json:"directory" yaml:"directory"`
	Algorithm string    `json:"algorithm" yaml:"algorithm"`
}

type File struct {
	Name   string `json:"file" yaml:"file"`
	Path   string `json:"-" yaml:"-"`
	Hash   string `json:"hash" yaml:"hash"`
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

type Files struct {
	Split []*File `json:"split,omitempty" yaml:"split,omitempty"`
	Core  []*File `json:"core,omitempty" yaml:"core,omitempty"`
}

type Output struct {
	Files []*File `json:"files" yaml:"files"`
}

type Integrity struct {
	Version  int      `json:"version" yaml:"version"`
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Files    Files    `json:"files" yaml:"files"`

	expectedFiles []*File
	validateFiles []*File
	combinedFiles []*File

	ignore           []string
	validate         bool
	directory        string
	filename         string
	getFirstExists   bool
	getFirstEmpty    bool
	getFirstValidate bool
}

func New(directory string, validate bool) (*Integrity, error) {
	abs, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}

	i := &Integrity{
		Version: 1,
		Metadata: Metadata{
			CreatedAt: time.Now().UTC(),
			Version:   common.AppVersion.Version,
			Directory: filepath.ToSlash(abs),
		},

		ignore:   []string{},
		validate: validate,

		getFirstExists:   false,
		getFirstEmpty:    false,
		getFirstValidate: false,
	}

	i.directory = filepath.ToSlash(directory)
	i.filename = filepath.Join(directory, common.Filename)

	if i.validate {
		if _, err := os.Stat(i.filename); err == nil {
			if err := i.LoadFile(); err != nil {
				return i, err
			}
		}
	}

	getFirstPath := filepath.Join(i.directory, common.GetFirstDirectory)

	if _, err := os.Stat(getFirstPath); err == nil {
		isEmpty, err := utils.IsDirectoryEmpty(getFirstPath)
		if err != nil {
			return i, err
		}

		i.getFirstExists = true
		i.getFirstEmpty = isEmpty
	}

	if i.validate && i.getFirstExists && !i.getFirstEmpty {
		i.getFirstValidate = true
	}

	i.expectedFiles = append(i.Files.Split, i.Files.Core...)

	return i, nil
}

func (i *Integrity) SetAlgorithm(algorithm string) bool {
	i.Metadata.Algorithm = algorithm
	return true
}

func (i *Integrity) Checks() error {
	if i.getFirstExists && i.getFirstEmpty {
		return fmt.Errorf("%s exists and is empty, this is not allowed, please delete or populate files", common.GetFirstDirectory)
	}

	return nil
}

func (i *Integrity) LoadFile() error {
	yamlFile, err := ioutil.ReadFile(i.filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlFile, i); err != nil {
		return err
	}

	return nil
}

func (i *Integrity) SortFiles() error {
	for _, f := range i.expectedFiles {
		if strings.HasPrefix(f.Name, common.GetFirstDirectory) {
			i.Files.Split = append(i.Files.Split, f)
		} else {
			i.Files.Core = append(i.Files.Core, f)
		}
	}

	return nil
}

func (i *Integrity) WriteFile() error {
	i.SortFiles()

	data, err := yaml.Marshal(i)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(i.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

func (i *Integrity) CompareFiles() (identical bool, err error) {
	identical = true
	skippedSplit := false

	expected := map[string]*File{}
	actual := map[string]*File{}
	combined := map[string]*File{}

	for _, file := range i.validateFiles {
		actual[file.Name] = file
	}

	for _, file := range i.expectedFiles {
		expected[file.Name] = file
	}

	for _, file := range i.validateFiles {
		if _, ok := expected[file.Name]; !ok {
			logrus.WithField("file", file.Name).WithField("status", "added").Warn("Added File")

			file.Status = "added"
			combined[file.Name] = file
			identical = false
		}
	}

	for _, file := range i.expectedFiles {
		if _, ok := actual[file.Name]; !ok {
			if strings.HasPrefix(file.Name, common.GetFirstDirectory) && !i.getFirstValidate {
				skippedSplit = true
				logrus.WithField("file", file.Name).Debugf("skipping split file because directory %s does not exist", common.GetFirstDirectory)
				continue
			}

			logrus.WithField("status", "missing").WithField("file", file.Name).Warn("Missing File")
			file.Status = "missing"
			combined[file.Name] = file
			identical = false
		}
	}

	for _, file := range i.validateFiles {
		if ef, ok := expected[file.Name]; ok {
			if ef.Hash != file.Hash {
				logrus.
					WithField("file", file.Name).
					WithField("status", "mismatch").
					WithField("expected_hash", ef.Hash).
					WithField("actual_hash", file.Hash).
					Error("Checksum Failure")

				file.Status = "failed"
				combined[file.Name] = file
				identical = false
			} else {
				logrus.
					WithField("file", file.Name).
					WithField("status", "ok").
					WithField("hash", ef.Hash).
					Debug("Checksum Validated")

				file.Status = "ok"
				combined[file.Name] = file
			}
		}
	}

	for _, f := range combined {
		i.combinedFiles = append(i.combinedFiles, f)
	}

	if skippedSplit {
		logrus.Warnf("Split files skipped as %s was missing or empty", common.GetFirstDirectory)
	}

	return identical, nil
}

func (i *Integrity) HashFiles() error {
	var files []*File

	if i.validate {
		files = i.validateFiles
	} else {
		files = i.expectedFiles
	}

	for _, file := range files {
		log := logrus.WithField("file", file.Name)

		hash, err := utils.HashFileSha256(file.Path)
		if err != nil {
			log.WithError(err).Error("unable to hash file")
			return err
		}

		log.WithField("hash", hash).Infof("Processed File")

		file.Hash = hash
	}

	return nil
}

func (i *Integrity) DiscoverFiles() error {
	var err error

	if i.validate {
		i.validateFiles, err = i.getFiles()
		if err != nil {
			return err
		}
	} else {
		i.expectedFiles, err = i.getFiles()
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Integrity) getFiles() (files []*File, err error) {
	if err := filepath.Walk(i.directory,
		func(path string, info os.FileInfo, err error) error {
			pathCheck, err := os.Stat(path)
			if err != nil {
				return err
			}

			if !pathCheck.IsDir() {
				fileName, err := filepath.Rel(i.directory, path)
				if err != nil {
					return err
				}

				for _, ignore := range i.ignore {
					if fileName == ignore {
						return nil
					}
					if strings.HasPrefix(fileName, ignore) {
						return nil
					}
					if matched, _ := regexp.MatchString(ignore, fileName); matched {
						return nil
					}
				}

				// Both name and path must be ToSlash because the Name is what
				// is ultimately written to the versioning file
				files = append(files, &File{
					Name: filepath.ToSlash(fileName),
					Path: filepath.ToSlash(path),
				})
			}

			return nil
		},
	); err != nil {
		return nil, err
	}

	return files, nil
}

func (i *Integrity) SetName(name string) error {
	match, err := regexp.MatchString("[0-9]{3}.[0-9]{2}.[0-9][A-Z]?", name)
	if err != nil {
		return err
	}

	if !match {
		return fmt.Errorf("%s does not match the required format", name)
	}

	i.Metadata.Name = name

	return nil
}

func (i *Integrity) SetUser(user string) {
	i.Metadata.CreatedBy = user
}

func (i *Integrity) SetIgnore(ignore []string) {
	i.ignore = ignore

	i.ignore = append(i.ignore, common.IgnoreAlways...)
}

func (i *Integrity) GetValidationOutput(format string) ([]byte, error) {
	if format != "json" {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if format == "json" {
		b, err := json.Marshal(Output{Files: i.combinedFiles})
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	return nil, nil
}
