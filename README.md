# ZKits Configurator Library #

[![ZKits](https://img.shields.io/badge/ZKits-Library-f3c)](https://github.com/edoger/zkits-configurator)
[![Build Status](https://travis-ci.org/edoger/zkits-configurator.svg?branch=master)](https://travis-ci.org/edoger/zkits-configurator)
[![Coverage Status](https://coveralls.io/repos/github/edoger/zkits-configurator/badge.svg?branch=master)](https://coveralls.io/github/edoger/zkits-configurator?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/11e8102293d44ede913f7f47603210ef)](https://www.codacy.com/manual/edoger/zkits-configurator?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=edoger/zkits-configurator&amp;utm_campaign=Badge_Grade)
[![Golang Version](https://img.shields.io/badge/golang-1.13+-orange)](https://github.com/edoger/zkits-configurator)

## About ##

This package is a library of ZKits project.
This library provides an efficient configuration file loading process controller. 
Built-in config file binding that supports TOML/JSON/XML/YAML format.

## Usage ##

 1. Install package.
 
    ```sh
    go get -u -v github.com/edoger/zkits-configurator
    ```
 
 2. Load configuration target.
 
    ```go
    // Create a new configurator.
    c := configurator.New()

    // Add configuration files.
    err := c.AddFile("/path/to/*.json")
    err := c.AddFile("/path/to/*.yaml")
    err := c.AddFile("/path/to/*.toml")
    err := c.AddFile("/path/to/*.xml")

    // Load configuration file.
    err := c.LoadJSON("file.name", object)
    err := c.LoadXML("file.name", object)
    err := c.LoadTOML("file.name", object)
    err := c.LoadYAML("file.name", object)
    
    item, err := c.Load("file.name")
    item.Bytes()
    item.String()
    item.JSON(object)
    item.XML(object)
    item.TOML(object)
    item.YAML(object)
    ```

## License ##

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)
