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
)

type FileLoader struct {
	dir string // The configuration file storage directory.
	ext string // The configuration file extension name.
}

// Create a new file loader.
// Requires that you specify the directory and file extension for the file search.
func NewFileLoader(dir, ext string) Loader {
	return &FileLoader{dir: dir, ext: ext}
}

func (loader *FileLoader) Load(target string, next func() ([]byte, error)) ([]byte, error) {
	var file string

	if loader.ext == "" {
		file = filepath.Join(loader.dir, target)
	} else {
		file = filepath.Join(loader.dir, target+"."+loader.ext)
	}

	if data, err := ioutil.ReadFile(file); err == nil {
		return data, nil
	} else {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return next()
	}
}
