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
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	o := New()
	if o == nil {
		t.Fatal("New(): nil")
	}
}

func TestConfigurator(t *testing.T) {
	o := New()
	l := NewFileLoader()
	if err := o.AddFile("test/*"); err != nil {
		t.Fatal(err)
	}
	o.Use(l)
	do := func(fs ...func(Configurator)) {
		for _, f := range fs {
			f(o)
		}
	}

	type Value struct {
		Name string `json:"name" xml:"name" toml:"name" yaml:"name"`
	}

	do(func(c Configurator) {
		item, err := c.Load("test.txt")
		if err != nil {
			t.Fatalf("Configurator.Load(): %s", err)
		}
		if got := item.String(); got != "test" {
			t.Fatalf("Configurator.Load(): %s", got)
		}

		if item, err := c.Load("unknown"); err == nil {
			t.Fatalf("Configurator.Load(): no error")
		} else {
			if err != ErrNotFound {
				t.Fatalf("Configurator.Load(): %s", err)
			}
			if item != nil {
				t.Fatalf("Configurator.Load(): %v", item)
			}
		}

	}, func(c Configurator) {
		v := new(Value)
		if err := c.LoadJSON("test.json", v); err != nil {
			t.Fatalf("Configurator.LoadJSON(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Configurator.LoadJSON(): %s", v.Name)
		}

		v = new(Value)
		if err := c.LoadJSON("unknown.json", v); err == nil {
			t.Fatalf("Configurator.LoadJSON(): no error")
		} else {
			if err != ErrNotFound {
				t.Fatalf("Configurator.LoadJSON(): %s", err)
			}
		}
		if v.Name != "" {
			t.Fatalf("Configurator.LoadJSON(): %s", v.Name)
		}
	}, func(c Configurator) {
		v := new(Value)
		if err := c.LoadXML("test.xml", v); err != nil {
			t.Fatalf("Configurator.LoadXML(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Configurator.LoadXML(): %s", v.Name)
		}

		v = new(Value)
		if err := c.LoadXML("unknown.xml", v); err == nil {
			t.Fatalf("Configurator.LoadXML(): no error")
		} else {
			if err != ErrNotFound {
				t.Fatalf("Configurator.LoadXML(): %s", err)
			}
		}
		if v.Name != "" {
			t.Fatalf("Configurator.LoadXML(): %s", v.Name)
		}
	}, func(c Configurator) {
		v := new(Value)
		if err := c.LoadTOML("test.toml", v); err != nil {
			t.Fatalf("Configurator.LoadTOML(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Configurator.LoadTOML(): %s", v.Name)
		}

		v = new(Value)
		if err := c.LoadTOML("unknown.toml", v); err == nil {
			t.Fatalf("Configurator.LoadTOML(): no error")
		} else {
			if err != ErrNotFound {
				t.Fatalf("Configurator.LoadTOML(): %s", err)
			}
		}
		if v.Name != "" {
			t.Fatalf("Configurator.LoadTOML(): %s", v.Name)
		}
	}, func(c Configurator) {
		v := new(Value)
		if err := c.LoadYAML("test.yaml", v); err != nil {
			t.Fatalf("Configurator.LoadYAML(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Configurator.LoadYAML(): %s", v.Name)
		}

		v = new(Value)
		if err := c.LoadYAML("unknown.yaml", v); err == nil {
			t.Fatalf("Configurator.LoadYAML(): no error")
		} else {
			if err != ErrNotFound {
				t.Fatalf("Configurator.LoadYAML(): %s", err)
			}
		}
		if v.Name != "" {
			t.Fatalf("Configurator.LoadYAML(): %s", v.Name)
		}
	})
}

func TestConfiguratorLoadError(t *testing.T) {
	o := New()
	o.Use(LoaderFunc(func(string) (Item, error) {
		return nil, errors.New("test")
	}))

	item, err := o.Load("test")
	if err == nil {
		t.Fatal("nil error")
	}
	if item != nil {
		t.Fatal(item)
	}
}

func TestConfiguratorLoad(t *testing.T) {
	o := New()
	if err := o.AddFile("test/foo/*"); err != nil {
		t.Fatal(err)
	}

	type Value struct {
		Name string `json:"name"`
	}

	do := func(fs ...func()) {
		for _, f := range fs {
			f()
		}
	}

	do(func() {
		v := new(Value)
		if err := o.LoadJSON("a", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "a.json" {
			t.Fatal(v.Name)
		}
	}, func() {
		v := new(Value)
		if err := o.LoadJSON("a.json", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "a.json" {
			t.Fatal(v.Name)
		}
	}, func() {
		v := new(Value)
		if err := o.LoadJSON("b", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "b.json" {
			t.Fatal(v.Name)
		}
	}, func() {
		v := new(Value)
		if err := o.LoadJSON("b.json", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "b.json" {
			t.Fatal(v.Name)
		}
	})

	o.Use(NewFileLoader().MustAddFile("test/bar/*"))

	do(func() {
		v := new(Value)
		if err := o.LoadJSON("a", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "bar" {
			t.Fatal(v.Name)
		}
	}, func() {
		v := new(Value)
		if err := o.LoadJSON("a.json", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "bar" {
			t.Fatal(v.Name)
		}
	}, func() {
		v := new(Value)
		if err := o.LoadJSON("b", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "b.json" {
			t.Fatal(v.Name)
		}
	}, func() {
		v := new(Value)
		if err := o.LoadJSON("b.json", v); err != nil {
			t.Fatal(err)
		}
		if v.Name != "b.json" {
			t.Fatal(v.Name)
		}
	})

	if _, err := o.Load("a.xml"); err != ErrNotFound {
		t.Fatal(err)
	}
}
