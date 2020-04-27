# ZKits Configurator Library #

[![Build Status](https://travis-ci.org/edoger/zkits-configurator.svg?branch=master)](https://travis-ci.org/edoger/zkits-configurator)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-configurator/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-configurator?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/11e8102293d44ede913f7f47603210ef)](https://www.codacy.com/manual/edoger/zkits-configurator?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-configurator&amp;utm_campaign=Badge_Grade)

## About ##

This package is a library of ZKits project.
This library provides an efficient configuration file loading process controller. 
Built-in configuration file binding that supports TOML / JSON / XML format.

## Usage ##

 1. Install package.
 
    ```sh
    go get -u -v github.com/edoger/zkits-configurator
    ```

 2. Register configuration loader.
 
    ```go
    // Create a new configuration manager client.
    client := configurator.New()
    // Create a configuration file loader.
    loader := configurator.NewFileLoader()

    // Add configuration directory.
    err := loader.AddDir("dir")
    err := loader.AddDir("dir", ".json", ".toml")
    err := loader.AddFile("path/to/file.ext")

    // Register the configuration loader.
    client.Use(loader)

    // Create memory configuration loader.
    memory := configurator.NewMemoryLoader()
    memory.Add("target", []byte(`{"key":"value"}`))
    memory.Set("target", []byte(`{"key":"value"}`))

    // Register the configuration loader.
    client.Use(loader)
    ```
    
 3. Load configuration target.

    ```go
    // Load the configuration target with the given name.
    content, err := client.Load("target")
    
    // Load the configuration target and binds to the given bindable.
    err := client.LoadAndBind("target", bindable)
    
    // Load a configuration file and return a *bytes.Buffer.
    buffer, err := client.LoadBuffer("file")
    
    // Load a JSON configuration and bind it to the given object.
    err := client.LoadJSON("file", object)
    
    // Load a XML configuration and bind it to the given object.
    err := client.LoadXML("file", object)
    
    // Load a TOML configuration and bind it to the given object.
    err := client.LoadTOML("file", object)
    ```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
