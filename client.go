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
	"sync"

	"github.com/BurntSushi/toml"
)

// ErrNotFound indicates that the configuration target was not found.
// When a configuration target cannot be loaded in all loaders, this error will be returned.
var ErrNotFound = errors.New("configuration not found")

// Client type is the configuration manager.
type Client struct {
	mutex   sync.RWMutex
	loaders []Loader
}

// New creates and returns a new client instance.
func New() *Client { return new(Client) }

// Next type defines the trigger for the loader.
type Next func() ([]byte, error)

// Loader interface type defines a configuration target loader.
// The loader is responsible for loading the latest content of a given configuration target.
type Loader interface {
	// The Load method loads the contents of a given configuration target.
	// If the target is not within the processing range, next should be called to give
	// the load permission to other configuration loaders.
	Load(target string, next Next) ([]byte, error)
}

// LoaderFunc type defines the configuration target loader function.
type LoaderFunc func(string, Next) ([]byte, error)

// Load method is used to call the configuration target loader function.
// This is just to implement the loader interface.
func (f LoaderFunc) Load(target string, next Next) ([]byte, error) {
	return f(target, next)
}

// Use method registers a configuration target loader.
func (c *Client) Use(loader Loader) *Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.loaders = append(c.loaders, loader)
	return c
}

// Load method will immediately load the given configuration target.
func (c *Client) Load(target string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.load(target, 0)
}

func (c *Client) load(target string, index int) ([]byte, error) {
	if n := index + 1; index >= len(c.loaders) {
		return nil, ErrNotFound
	} else {
		return c.loaders[index].Load(target, func() ([]byte, error) {
			return c.load(target, n)
		})
	}
}

// LoadBuffer method loads the given configuration target and returns the content as
// a *bytes.Buffer instance.
func (c *Client) LoadBuffer(target string) (*bytes.Buffer, error) {
	if data, err := c.Load(target); err != nil {
		return nil, err
	} else {
		return bytes.NewBuffer(data), nil
	}
}

// LoadJSON method loads the given configuration target and binds the content to the
// given object as json.
func (c *Client) LoadJSON(target string, o interface{}) error {
	if data, err := c.Load(target); err != nil {
		return err
	} else {
		return json.Unmarshal(data, o)
	}
}

// LoadXML method loads the given configuration target and binds the content to the
// given object as xml.
func (c *Client) LoadXML(target string, o interface{}) error {
	if data, err := c.Load(target); err != nil {
		return err
	} else {
		return xml.Unmarshal(data, o)
	}
}

// LoadTOML method loads the given configuration target and binds the content to the
// given object as toml.
func (c *Client) LoadTOML(target string, o interface{}) error {
	if data, err := c.Load(target); err != nil {
		return err
	} else {
		return toml.Unmarshal(data, o)
	}
}
