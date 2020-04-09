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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var _ Loader = (*FileLoader)(nil)

// FileLoader type is a built-in configuration file loader.
type FileLoader struct {
	// Resource read-write lock.
	mutex sync.RWMutex
	// The path of the configuration file directory.
	directory string
	// Whether the loader has been initialized.
	initialized bool
	// Configuration file path cache.
	// This is a mapping of names (without extension name) to paths.
	files map[string][]string
}

// Create an instance of the configuration file loader.
func NewFileLoader(directory string) *FileLoader {
	return &FileLoader{directory: directory}
}

// Initialize the current configuration file loader.
// This method will build a path cache table of all regular configuration
// files in the directory.
// This method is idempotent.
func (loader *FileLoader) Initialize() error {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	if loader.initialized {
		return nil
	}

	dir, err := filepath.Abs(loader.directory)
	if err != nil {
		return err
	}

	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if loader.files == nil {
		loader.files = make(map[string][]string, len(items))
	}
	for i, j := 0, len(items); i < j; i++ {
		if items[i].Mode().IsRegular() {
			// We only care about regular file.
			name := items[i].Name()
			if ext := filepath.Ext(name); ext != "" {
				name = strings.TrimSuffix(name, ext)
			}
			loader.files[name] = append(loader.files[name], filepath.Join(dir, items[i].Name()))
		}
	}

	loader.initialized = true
	return nil
}

// Load a configuration file by the given file name.
func (loader *FileLoader) Load(target string, next func() ([]byte, error)) ([]byte, error) {
	loader.mutex.RLock()
	defer loader.mutex.RUnlock()

	if len(loader.files) > 0 {
		for _, item := range loader.files[target] {
			data, err := ioutil.ReadFile(item)
			if err == nil {
				return data, nil
			}
			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}

	return next()
}
