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

	"github.com/BurntSushi/toml"
)

// Client type is the configuration file loading manager.
type Client struct {
	// List of registered configuration loaders.
	loaders []Loader
}

// Create a new configuration manager client.
func New() *Client {
	return new(Client)
}

// Register a configuration loader.
// List of registered configuration item loaders.
func (c *Client) Use(loader Loader) *Client {
	c.loaders = append(c.loaders, loader)
	return c
}

// Load the configuration item with the given name.
func (c *Client) Load(target string) ([]byte, error) {
	return c.load(target, 0)
}

func (c *Client) load(target string, index int) ([]byte, error) {
	if next := index + 1; index >= len(c.loaders) {
		return nil, ErrNotFound
	} else {
		return c.loaders[index].Load(target, func() ([]byte, error) {
			return c.load(target, next)
		})
	}
}

// Load a configuration file and return a *bytes.Buffer.
// If the load fails, a nil will be returned.
func (c *Client) LoadBuffer(target string) (*bytes.Buffer, error) {
	if data, err := c.Load(target); err != nil {
		return nil, err
	} else {
		return bytes.NewBuffer(data), nil
	}
}

// Load a JSON configuration and bind it to the given object.
func (c *Client) LoadJSON(target string, o interface{}) error {
	if data, err := c.Load(target); err != nil {
		return err
	} else {
		return json.Unmarshal(data, o)
	}
}

// Load a XML configuration and bind it to the given object.
func (c *Client) LoadXML(target string, o interface{}) error {
	if data, err := c.Load(target); err != nil {
		return err
	} else {
		return xml.Unmarshal(data, o)
	}
}

// Load a TOML configuration and bind it to the given object.
func (c *Client) LoadTOML(target string, o interface{}) error {
	if data, err := c.Load(target); err != nil {
		return err
	} else {
		return toml.Unmarshal(data, o)
	}
}
