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
	"errors"
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
