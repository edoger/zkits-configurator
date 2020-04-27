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
	"sync"
)

// Make sure that the memory loader implements the loader interface.
var _ Loader = (*MemoryLoader)(nil)

// MemoryLoader type defines a memory configuration loader.
type MemoryLoader struct {
	mutex sync.RWMutex
	items map[string][]byte
}

// NewMemoryLoader creates and returns a new memory configuration loader instance.
func NewMemoryLoader() *MemoryLoader {
	return &MemoryLoader{items: make(map[string][]byte)}
}

// Add method adds a given configuration data to the current loader.
// If the given configuration target already exists, it will not be added.
// This method returns whether the given configuration was added successfully.
func (loader *MemoryLoader) Add(target string, value []byte) bool {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	if data, found := loader.items[target]; found {
		return bytes.Equal(data, value)
	} else {
		loader.set(target, value)
		return true
	}
}

// Set method adds a given configuration data to the current loader.
// If the given configuration item already exists, the old configuration data is overwritten.
// If the given configuration data is nil, delete the corresponding configuration data.
func (loader *MemoryLoader) Set(target string, value []byte) {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	loader.set(target, value)
}

// Set the given configuration data.
func (loader *MemoryLoader) set(target string, value []byte) {
	if value == nil {
		delete(loader.items, target)
	} else {
		loader.items[target] = make([]byte, len(value))
		copy(loader.items[target], value)
	}
}

// Load method loads the latest content of a given configuration target.
func (loader *MemoryLoader) Load(target string, next Next) ([]byte, error) {
	if data := loader.load(target); data != nil {
		return data, nil
	}
	return next()
}

// Search and load the given configuration target content.
func (loader *MemoryLoader) load(target string) []byte {
	loader.mutex.RLock()
	defer loader.mutex.RUnlock()

	if data, found := loader.items[target]; found {
		r := make([]byte, len(data))
		copy(r, data)
		return r
	}
	return nil
}
