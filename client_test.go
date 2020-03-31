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
	"testing"
)

func TestClient(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("New() return nil")
	}

	l1 := LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "foo" {
			return []byte("foo"), nil
		}
		return next()
	})
	l2 := LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "bar" {
			return []byte("bar"), nil
		}
		return next()
	})
	if client.Use(l1) == nil {
		t.Fatal("Client.Use() return nil")
	}
	if client.Use(l2) == nil {
		t.Fatal("Client.Use() return nil")
	}

	if data, err := client.Load("test"); err == nil {
		t.Fatal("Client.Load()")
	} else {
		if err != ErrNotFound {
			t.Fatal(err)
		}
		if data != nil {
			t.Fatal(string(data))
		}
	}
	if data, err := client.Load("foo"); err != nil {
		t.Fatal(err)
	} else {
		if s := string(data); s != "foo" {
			t.Fatal(s)
		}
	}
	if data, err := client.Load("bar"); err != nil {
		t.Fatal(err)
	} else {
		if s := string(data); s != "bar" {
			t.Fatal(s)
		}
	}
}

func TestClient_LoadJSON(t *testing.T) {
	type object struct {
		Value string `json:"key"`
	}

	client := New()
	client.Use(LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "foo" {
			return []byte(`{"key":"foo"}`), nil
		}
		return next()
	}))

	o := new(object)

	if err := client.LoadJSON("foo", o); err != nil {
		t.Fatal(err)
	}
	if o.Value != "foo" {
		t.Fatal(o.Value)
	}

	o = new(object)
	if err := client.LoadJSON("bar", o); err == nil {
		t.Fatal("Client.LoadJSON()")
	} else {
		if err != ErrNotFound {
			t.Fatal(err)
		}
	}
	if o.Value != "" {
		t.Fatal(o.Value)
	}

	client.Use(LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "bar" {
			return []byte(`{bar}`), nil
		}
		return next()
	}))

	o = new(object)
	if err := client.LoadJSON("bar", o); err == nil {
		t.Fatal("Client.LoadJSON()")
	}
	if o.Value != "" {
		t.Fatal(o.Value)
	}
}

func TestClient_LoadXML(t *testing.T) {
	type object struct {
		Value string `xml:"key"`
	}

	client := New()
	client.Use(LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "foo" {
			return []byte(`<xml><key>foo</key></xml>`), nil
		}
		return next()
	}))

	o := new(object)

	if err := client.LoadXML("foo", o); err != nil {
		t.Fatal(err)
	}
	if o.Value != "foo" {
		t.Fatal(o.Value)
	}

	o = new(object)
	if err := client.LoadXML("bar", o); err == nil {
		t.Fatal("Client.LoadXML()")
	} else {
		if err != ErrNotFound {
			t.Fatal(err)
		}
	}
	if o.Value != "" {
		t.Fatal(o.Value)
	}

	client.Use(LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "bar" {
			return []byte(`><`), nil
		}
		return next()
	}))

	o = new(object)
	if err := client.LoadXML("bar", o); err == nil {
		t.Fatal("Client.LoadXML()")
	}
	if o.Value != "" {
		t.Fatal(o.Value)
	}
}

func TestClient_LoadTOML(t *testing.T) {
	type object struct {
		Value string `toml:"key"`
	}

	client := New()
	client.Use(LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "foo" {
			return []byte(`key = "foo"`), nil
		}
		return next()
	}))

	o := new(object)

	if err := client.LoadTOML("foo", o); err != nil {
		t.Fatal(err)
	}
	if o.Value != "foo" {
		t.Fatal(o.Value)
	}

	o = new(object)
	if err := client.LoadTOML("bar", o); err == nil {
		t.Fatal("Client.LoadTOML()")
	} else {
		if err != ErrNotFound {
			t.Fatal(err)
		}
	}
	if o.Value != "" {
		t.Fatal(o.Value)
	}

	client.Use(LoaderFunc(func(target string, next func() ([]byte, error)) ([]byte, error) {
		if target == "bar" {
			return []byte(`===`), nil
		}
		return next()
	}))

	o = new(object)
	if err := client.LoadTOML("bar", o); err == nil {
		t.Fatal("Client.LoadTOML()")
	}
	if o.Value != "" {
		t.Fatal(o.Value)
	}
}
