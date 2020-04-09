# ZKits Configurator Library #

[![Build Status](https://travis-ci.org/edoger/zkits-configurator.svg?branch=master)](https://travis-ci.org/edoger/zkits-configurator)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-configurator/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-configurator?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/11e8102293d44ede913f7f47603210ef)](https://www.codacy.com/manual/edoger/zkits-configurator?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-configurator&amp;utm_campaign=Badge_Grade)

## About ##

This package is a library of ZKits project.
This library provides an efficient configuration file loading process controller. 
Built-in configuration file binding that supports TOML / JSON / XML format.

## Usage ##

 1. Import package.
 
    ```sh
    go get -u -v github.com/edoger/zkits-configurator
    ```

 2. Example.
    ```go
    package main
    
    import (
        "github.com/edoger/zkits-configurator"
    )
    
    func main() {
        // Create a new configuration manager client.
        client := configurator.New()
    
        // Create a configuration file loader.
        // We can customize the loader of various configuration data sources,
        // only need to implement the Loader interface.
        loader := configurator.NewFileLoader("config")
        // Initialize the configuration file loader.
        if err := loader.Initialize(); err != nil {
            // Handle error.
        }
    
        // Register the configuration loader.
        // The loader that is registered first has the highest priority.
        client.Use(loader)
    
        // Load the configuration target with the given name.
        content, err := client.Load("file-name")
        if err != nil {
            // Handle error.
        }
    
        // Load a configuration file and return a *bytes.Buffer.
        client.LoadBuffer("file")
    
        // Load a JSON configuration and bind it to the given object.
        client.LoadJSON("file", object)
        // Load a XML configuration and bind it to the given object.
        client.LoadXML("file", object)
        // Load a TOML configuration and bind it to the given object.
        client.LoadTOML("file", object)
    }
    ```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
