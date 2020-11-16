# ZKits Configurator Library #

[![ZKits](https://img.shields.io/badge/ZKits-Library-f3c)](https://github.com/edoger/zkits-configurator)
[![Build Status](https://travis-ci.org/edoger/zkits-configurator.svg?branch=master)](https://travis-ci.org/edoger/zkits-configurator)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-configurator/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-configurator?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/11e8102293d44ede913f7f47603210ef)](https://www.codacy.com/manual/edoger/zkits-configurator?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-configurator&amp;utm_campaign=Badge_Grade)
[![Golang Version](https://img.shields.io/badge/golang-1.13+-orange)](https://github.com/edoger/zkits-configurator)

## About ##

This package is a library of ZKits project. 
The configurator supports loading config targets from anywhere through built-in or custom loaders.

## Install ##

```sh
go get -u -v github.com/edoger/zkits-configurator
```

## Usage ##

```go
package main

import (
	"github.com/edoger/zkits-configurator"
)

func main() {
	// Create a new configurator.
	c := configurator.New()

	// Add config files.
	// The parameter format must be the format required by filepath.Glob().
	c.AddFile("/path/to/*.json")
	c.AddFile("/path/to/*.yaml")
	c.AddFile("/path/to/*.toml")
	c.AddFile("/path/to/*.xml")

	// Load config file by name.
	// Usually, the file ext name can be omitted, and the configurator is intelligent enough.
	item, err := c.Load("file.name")
	if err != nil {
		panic(err)
	}

	item.IsEmpty() // Is empty file?
	item.Len()     // File content length.
	item.Bytes()   // Get file content as []byte.
	item.String()  // Get file content as string.
	item.Reader()  // Get file content as io.Reader.

	// Variable used to bind config content.
	object := make(map[string]interface{})

	// Bind config content to object as json/xml/toml/yaml format.
	item.JSON(&object)
	item.XML(&object)
	item.TOML(&object)
	item.YAML(&object)
	// These methods can also be used!
	c.LoadJSON("file.name", &object)
	c.LoadXML("file.name", &object)
	c.LoadTOML("file.name", &object)
	c.LoadYAML("file.name", &object)

	// Register a custom configuration loader (which takes precedence over the built-in file loader).
	c.Use(configurator.LoaderFunc(func(target string) (configurator.Item, error) {
		// Do something!
		return nil, nil
	}))
	// We provide a built-in configuration file loader.
	loader := configurator.NewFileLoader()
	loader.MustAddFile("/path/to/other/*.json")
	// You can register more loaders, they just need to implement the configurator.Loader interface.
	c.Use(loader)
}
```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
