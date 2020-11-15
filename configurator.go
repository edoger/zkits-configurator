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

// Configurator defines the configuration manager.
type Configurator interface {
	// Use registers a custom configuration loader.
	// The last registered configuration loader will have the highest priority.
	Use(Loader) Configurator

	// AddFile adds one or more config files to the current loader.
	// The given parameter need to comply with the search rules supported
	// by filepath.Glob.
	// This method comes from the built-in configuration file loader.
	AddFile(string) error

	// Load loads the given config target.
	// If the given config target does not exist, ErrNotFound is returned.
	// We will give priority to the custom loader. If there is no available config loader
	// or all registered config loaders cannot load the config target (the semantic target
	// does not exist), it will be automatically delegated to the built-in config file loader.
	Load(string) (Item, error)

	// LoadJSON loads the given config target and binds it to the given object as json.
	LoadJSON(string, interface{}) error

	// LoadXML loads the given config target and binds it to the given object as xml.
	LoadXML(string, interface{}) error

	// LoadTOML loads the given config target and binds it to the given object as toml.
	LoadTOML(string, interface{}) error

	// LoadYAML loads the given config target and binds it to the given object as yaml.
	LoadYAML(string, interface{}) error
}

// New creates and returns a new Configurator instance.
func New() Configurator {
	fs := newFileLoader()
	return &configurator{fs: fs, loaders: []Loader{fs}}
}

// The configurator type is a built-in implementation of the Configurator interface.
type configurator struct {
	fs      *fileLoader
	loaders []Loader
}

// Use registers a custom configuration loader.
// The last registered configuration loader will have the highest priority.
func (o *configurator) Use(loader Loader) Configurator {
	o.loaders = append(o.loaders, loader)
	return o
}

// AddFile adds one or more config files to the current loader.
// The given parameter need to comply with the search rules supported
// by filepath.Glob.
// This method comes from the built-in configuration file loader.
func (o *configurator) AddFile(pattern string) error {
	return o.fs.AddFile(pattern)
}

// Load loads the given config target.
// If the given config target does not exist, ErrNotFound is returned.
// We will give priority to the custom loader. If there is no available config loader
// or all registered config loaders cannot load the config target (the semantic target
// does not exist), it will be automatically delegated to the built-in config file loader.
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

// LoadJSON loads the given config target and binds it to the given object as json.
func (o *configurator) LoadJSON(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.JSON(v)
	}
}

// LoadXML loads the given config target and binds it to the given object as xml.
func (o *configurator) LoadXML(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.XML(v)
	}
}

// LoadTOML loads the given config target and binds it to the given object as toml.
func (o *configurator) LoadTOML(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.TOML(v)
	}
}

// LoadYAML loads the given config target and binds it to the given object as yaml.
func (o *configurator) LoadYAML(target string, v interface{}) error {
	if item, err := o.Load(target); err != nil {
		return err
	} else {
		return item.YAML(v)
	}
}
