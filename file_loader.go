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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Make sure that the file loader implements the loader interface.
var _ Loader = (*FileLoader)(nil)

// FileLoader type defines a configuration file loader.
// The loader is used to search and load the latest content of the configuration
// file matching the target in the known file list.
type FileLoader struct {
	mutex sync.RWMutex
	files map[string]string
}

// NewFileLoader creates and returns a new file loader instance.
func NewFileLoader() *FileLoader {
	return &FileLoader{files: make(map[string]string)}
}

// AddDir method adds the given directory to the current file loader.
// This method can specify the extension names to limit the file type.
func (loader *FileLoader) AddDir(directory string, extensions ...string) error {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	dir, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	if files, err := loader.dir(dir, extensions); err != nil {
		return err
	} else {
		// Make sure that the file name added eventually will not be duplicated.
		for name, position := range files {
			if old, found := loader.files[name]; found {
				if old != position {
					return fmt.Errorf("duplicate file %s and %s", old, position)
				}
				delete(files, name)
			}
		}
		for name, position := range files {
			loader.files[name] = position
		}
		return nil
	}
}

// Read the list of regular files in the given directory.
func (loader *FileLoader) dir(directory string, extensions []string) (map[string]string, error) {
	items, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	files := make(map[string]string)
	for i, j := 0, len(items); i < j; i++ {
		// When adding directory, we only care about regular files.
		if mode := items[i].Mode(); !mode.IsRegular() {
			continue
		}
		n := items[i].Name()
		p := filepath.Join(directory, n)
		e := filepath.Ext(n)
		if e != "" {
			n = strings.TrimSuffix(n, e)
		}
		// Determines if the file extensions match.
		// If the extension is not qualified, all files will be added.
		if l := len(extensions); l > 0 {
			ok := false
			for i := 0; i < l; i++ {
				if extensions[i] == e {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		// Make sure that there are no duplicate files in the current directory.
		if old, found := files[n]; found {
			return nil, fmt.Errorf("duplicate file %s and %s", old, p)
		} else {
			files[n] = p
		}
	}

	return files, nil
}

// AddFile method adds a regular file to the current file loader.
func (loader *FileLoader) AddFile(name string) error {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	position, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	info, err := os.Stat(position)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("file %s is not a regular file", name)
	}

	n := info.Name()
	if e := filepath.Ext(n); e != "" {
		n = strings.TrimSuffix(n, e)
	}

	if old, found := loader.files[n]; found {
		if old != position {
			return fmt.Errorf("duplication file name %s and %s", old, position)
		}
	} else {
		loader.files[n] = position
	}
	return nil
}

// Load method loads the latest content of a given configuration file target.
func (loader *FileLoader) Load(target string, next Next) ([]byte, error) {
	content, err := loader.load(target)
	if err != nil {
		return nil, err
	}
	if content != nil {
		return content, nil
	}
	return next()
}

// Search and load the given file target content.
func (loader *FileLoader) load(target string) ([]byte, error) {
	loader.mutex.RLock()
	defer loader.mutex.RUnlock()

	if position, found := loader.files[target]; found {
		content, err := ioutil.ReadFile(position)
		if err == nil {
			return content, nil
		}
		// Be careful if the file is deleted?
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return nil, nil
}
