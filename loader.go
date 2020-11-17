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
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Loader interface defines the config target loader.
type Loader interface {
	// Load loads the given config target.
	// If this method returns a non-nil error, the entire configuration search and
	// loading process will be terminated immediately. If the returned Item is nil,
	// the next loader in the queue will be automatically run.
	Load(string) (Item, error)
}

// LoaderFunc type defines the config target function loader.
type LoaderFunc func(string) (Item, error)

// Load loads the given config target.
func (f LoaderFunc) Load(target string) (Item, error) {
	return f(target)
}

// FileLoader interface defines the config file loader.
type FileLoader interface {
	Loader

	// AddFile adds one or more config files to the current loader.
	// The given parameter need to comply with the search rules supported
	// by filepath.Glob.
	AddFile(string) error

	// MustAddFile adds one or more config files to the current loader.
	// This method is very similar to AddFile, the only difference is that it panics
	// when the add fails.
	MustAddFile(pattern string) FileLoader
}

// NewFileLoader creates and returns a config file loader instance.
func NewFileLoader() FileLoader {
	return newFileLoader()
}

// The newFileLoader function creates and returns a new fileLoader instance.
func newFileLoader() *fileLoader {
	return new(fileLoader)
}

// The fileLoader type is a built-in implementation of the FileLoader interface.
type fileLoader struct {
	mutex sync.RWMutex
	files map[string][][4]string
}

// AddFile adds one or more config files to the current loader.
// The given parameter need to comply with the search rules supported
// by filepath.Glob.
func (o *fileLoader) AddFile(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return err
	}

	var list [][4]string
	for i, j := 0, len(matches); i < j; i++ {
		info, err := os.Stat(matches[i])
		if err != nil {
			return err
		}
		// We only care about regular files.
		if info.Mode().IsRegular() {
			b := info.Name()
			e := filepath.Ext(b)
			list = append(list, [4]string{e, strings.TrimSuffix(b, e), b, matches[i]})
		}
	}

	if len(list) > 0 {
		o.mutex.Lock()
		if o.files == nil {
			o.files = make(map[string][][4]string)
		}
		for i, j := 0, len(list); i < j; i++ {
			o.files[list[i][1]] = append(o.files[list[i][1]], list[i])
		}
		o.mutex.Unlock()
	}
	return nil
}

// MustAddFile adds one or more config files to the current loader.
// This method is very similar to AddFile, the only difference is that it panics
// when the add fails.
func (o *fileLoader) MustAddFile(pattern string) FileLoader {
	if err := o.AddFile(pattern); err != nil {
		panic(err)
	}
	return o
}

// Load loads the given config file target.
// If the given config file does not exist, nil Item is returned.
func (o *fileLoader) Load(target string) (Item, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	if len(o.files) == 0 {
		return nil, nil
	}

	// If the config target already exists, the ext name will not be split and
	// the latest config file will be returned directly.
	// For example: Given "name.suffix", returns "/path/to/name.suffix.json".
	if a := o.files[target]; len(a) > 0 {
		item, err := newFileItem(a[len(a)-1][3])
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	e := filepath.Ext(target)
	// If there is no ext name, it can be determined that the target does not exist.
	if e == "" {
		return nil, nil
	}

	n := strings.TrimSuffix(target, e)
	if len(o.files[n]) == 0 {
		return nil, nil
	}

	r := o.files[n]
	for i := len(r) - 1; i >= 0; i-- {
		if r[i][0] == e {
			item, err := newFileItem(r[i][3])
			if err != nil {
				return nil, err
			}
			return item, nil
		}
	}
	return nil, nil
}
