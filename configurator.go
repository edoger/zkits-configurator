// Copyright 2020 The ZKits Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configurator

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

var (
	// ErrEmptyItem reports that the config item content is empty.
	// This error is returned when trying to bind an empty configuration item to any object.
	ErrEmptyItem = errors.New("configurator: empty item")

	// ErrNotFound indicates that the configuration target was not found.
	// When a configuration target cannot be loaded in all loaders, this error will be returned.
	ErrNotFound = errors.New("configurator: not found")
)

type Configurator interface {
	Use(Loader) Configurator

	AddFile(string) error

	Load(string) (Item, error)
	LoadJSON(string, interface{}) error
	LoadXML(string, interface{}) error
	LoadTOML(string, interface{}) error
	LoadYAML(string, interface{}) error
}

func New() Configurator {
	fs := newFileLoader()
	return &configurator{fs: fs, loaders: []Loader{fs}}
}

type configurator struct {
	fs *fileLoader

	loaders []Loader
}

func (o *configurator) Use(loader Loader) Configurator {
	o.loaders = append(o.loaders, loader)
	return o
}

func (o *configurator) AddFile(pattern string) error {
	return o.fs.AddFile(pattern)
}

func (o *configurator) Load(target string) (Item, error) {
	for k := len(o.loaders) - 1; k >= 0; k-- {
		if item, err := o.loaders[k].Load(target); err == nil {
			if item != nil {
				return item, nil
			}
		} else {
			if err != ErrNotFound {
				return nil, err
			}
		}
	}
	return nil, ErrNotFound
}

func (o *configurator) LoadJSON(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.JSON(v)
	}
}

func (o *configurator) LoadXML(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.XML(v)
	}
}

func (o *configurator) LoadTOML(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.TOML(v)
	}
}

func (o *configurator) LoadYAML(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.YAML(v)
	}
}

// Item interface defines the config item.
type Item interface {
	// IsEmpty determines whether the current config item content is empty.
	IsEmpty() bool

	// Len returns the current config item content length.
	Len() int

	// Bytes returns the current config item content bytes.
	Bytes() []byte

	// Bytes returns the current config item content string.
	String() string

	// Bytes returns the current config item content reader.
	Reader() io.Reader

	// JSON binds the current config item to the given object as json format.
	// If the current configuration item is empty, ErrEmptyItem will be returned.
	JSON(interface{}) error

	// XML binds the current config item to the given object as xml format.
	// If the current configuration item is empty, ErrEmptyItem will be returned.
	XML(interface{}) error

	// TOML binds the current config item to the given object as toml format.
	// If the current configuration item is empty, ErrEmptyItem will be returned.
	TOML(interface{}) error

	// YAML binds the current config item to the given object as yaml format.
	// If the current configuration item is empty, ErrEmptyItem will be returned.
	YAML(interface{}) error
}

// NewItemFromBytes creates and returns a config item from the given bytes.
func NewItemFromBytes(data []byte) Item {
	return newBytesItem(data)
}

// NewItemFromString creates and returns a config item from the given string.
func NewItemFromString(s string) Item {
	return newBytesItem([]byte(s))
}

// NewItemFromReader creates and returns a config item from the given reader.
func NewItemFromReader(r io.Reader) (Item, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return newBytesItem(data), nil
}

// The newBytesItem function creates and returns a config item from the given bytes.
func newBytesItem(data []byte) *bytesItem {
	return &bytesItem{data}
}

// The bytesItem type is a built-in implementation of the Item interface.
type bytesItem struct {
	data []byte
}

// IsEmpty determines whether the current config item content is empty.
func (item *bytesItem) IsEmpty() bool {
	return len(item.data) == 0
}

// Len returns the current config item content length.
func (item *bytesItem) Len() int {
	return len(item.data)
}

// Bytes returns the current config item content bytes.
func (item *bytesItem) Bytes() []byte {
	return item.data
}

// Bytes returns the current config item content string.
func (item *bytesItem) String() string {
	return string(item.data)
}

// Bytes returns the current config item content reader.
func (item *bytesItem) Reader() io.Reader {
	return bytes.NewReader(item.data)
}

// JSON binds the current config item to the given object as json format.
// If the current configuration item is empty, ErrEmptyItem will be returned.
func (item *bytesItem) JSON(o interface{}) error {
	if len(item.data) == 0 {
		return ErrEmptyItem
	}
	return json.Unmarshal(item.data, o)
}

// XML binds the current config item to the given object as xml format.
// If the current configuration item is empty, ErrEmptyItem will be returned.
func (item *bytesItem) XML(o interface{}) error {
	if len(item.data) == 0 {
		return ErrEmptyItem
	}
	return xml.Unmarshal(item.data, o)
}

// TOML binds the current config item to the given object as toml format.
// If the current configuration item is empty, ErrEmptyItem will be returned.
func (item *bytesItem) TOML(o interface{}) error {
	if len(item.data) == 0 {
		return ErrEmptyItem
	}
	return toml.Unmarshal(item.data, o)
}

// YAML binds the current config item to the given object as yaml format.
// If the current configuration item is empty, ErrEmptyItem will be returned.
func (item *bytesItem) YAML(o interface{}) error {
	if len(item.data) == 0 {
		return ErrEmptyItem
	}
	return yaml.Unmarshal(item.data, o)
}

// FileItem interface defines the config file item.
type FileItem interface {
	Item

	// Path returns the config file absolute path.
	Path() string

	// Base returns the config file base name.
	Base() string

	// Name returns the config file name (without extension name).
	Name() string
}

// NewFileItem returns the config file item created from the given file path.
// An error is returned only if a failure occurs to read the contents of the
// given config file.
func NewFileItem(s string) (FileItem, error) {
	path, err := filepath.Abs(s)
	if err != nil {
		return nil, err
	}
	return newFileItem(path)
}

// The newFileItem function reads the contents of the given file and returns
// a config file item.
func newFileItem(path string) (*fileItem, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	base := filepath.Base(path)
	return &fileItem{
		path,
		base,
		strings.TrimSuffix(base, filepath.Ext(base)),
		newBytesItem(data),
	}, nil
}

// The fileItem type is a built-in implementation of the FileItem interface.
type fileItem struct {
	path string
	base string
	name string

	*bytesItem
}

// Path returns the config file absolute path.
func (item *fileItem) Path() string {
	return item.path
}

// Base returns the config file base name.
func (item *fileItem) Base() string {
	return item.base
}

// Name returns the config file name (without extension name).
func (item *fileItem) Name() string {
	return item.name
}
