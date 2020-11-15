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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewItemFromBytes(t *testing.T) {
	o := NewItemFromBytes([]byte(""))
	if o == nil {
		t.Fatal("NewItemFromBytes(): nil")
	}
}

func TestNewItemFromString(t *testing.T) {
	o := NewItemFromString("")
	if o == nil {
		t.Fatal("NewItemFromString(): nil")
	}
}

func TestNewItemFromReader(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	o, err := NewItemFromReader(r)
	if err != nil {
		t.Fatalf("NewItemFromString(): %s", err)
	}
	if o == nil {
		t.Fatal("NewItemFromString(): nil")
	}
}

func TestBytesItem_IsEmpty(t *testing.T) {
	if got := NewItemFromString("").IsEmpty(); !got {
		t.Fatalf("Item.IsEmpty(): %v", got)
	}
	if got := NewItemFromString("test").IsEmpty(); got {
		t.Fatalf("Item.IsEmpty(): %v", got)
	}
}

func TestBytesItem_Len(t *testing.T) {
	if got := NewItemFromString("").Len(); got != 0 {
		t.Fatalf("Item.Len(): %v", got)
	}
	if got := NewItemFromString("test").Len(); got != 4 {
		t.Fatalf("Item.Len(): %v", got)
	}
}

func TestBytesItem_Bytes(t *testing.T) {
	if got := NewItemFromString("").Bytes(); !bytes.Equal(got, []byte("")) {
		t.Fatalf("Item.Bytes(): %v", got)
	}
	if got := NewItemFromString("test").Bytes(); !bytes.Equal(got, []byte("test")) {
		t.Fatalf("Item.Bytes(): %v", string(got))
	}
}

func TestBytesItem_String(t *testing.T) {
	if got := NewItemFromString("").String(); got != "" {
		t.Fatalf("Item.String(): %v", got)
	}
	if got := NewItemFromString("test").String(); got != "test" {
		t.Fatalf("Item.String(): %v", got)
	}
}

func TestBytesItem_Reader(t *testing.T) {
	r := NewItemFromString("test").Reader()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Item.Reader(): %s", err)
	}
	if !bytes.Equal(data, []byte("test")) {
		t.Fatalf("Item.Reader(): %v", string(data))
	}
}

func TestBytesItem(t *testing.T) {
	type Value struct {
		Name string `json:"name" xml:"name" toml:"name" yaml:"name"`
	}

	do := func(fs ...func(*Value)) {
		for _, f := range fs {
			f(new(Value))
		}
	}

	do(func(v *Value) {
		if err := NewItemFromString("").JSON(v); err != ErrEmptyItem {
			t.Fatalf("Item.JSON(): %s", err)
		}
		if err := NewItemFromString(`{"name":"test"}`).JSON(v); err != nil {
			t.Fatalf("Item.JSON(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Item.JSON(): %s", v.Name)
		}
	}, func(v *Value) {
		if err := NewItemFromString("").XML(v); err != ErrEmptyItem {
			t.Fatalf("Item.XML(): %s", err)
		}
		if err := NewItemFromString(`<xml><name>test</name></xml>`).XML(v); err != nil {
			t.Fatalf("Item.XML(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Item.XML(): %s", v.Name)
		}
	}, func(v *Value) {
		if err := NewItemFromString("").TOML(v); err != ErrEmptyItem {
			t.Fatalf("Item.TOML(): %s", err)
		}
		if err := NewItemFromString(`name = 'test'`).TOML(v); err != nil {
			t.Fatalf("Item.TOML(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Item.TOML(): %s", v.Name)
		}
	}, func(v *Value) {
		if err := NewItemFromString("").YAML(v); err != ErrEmptyItem {
			t.Fatalf("Item.YAML(): %s", err)
		}
		if err := NewItemFromString(`name: 'test'`).YAML(v); err != nil {
			t.Fatalf("Item.YAML(): %s", err)
		}
		if v.Name != "test" {
			t.Fatalf("Item.YAML(): %s", v.Name)
		}
	})
}

func TestNewFileItem(t *testing.T) {
	if item, err := NewFileItem(""); err == nil {
		t.Fatalf("NewFileItem(): no error")
	} else {
		if item != nil {
			t.Fatalf("NewFileItem(): %v", item)
		}
	}

	if item, err := NewFileItem("test/test.txt"); err != nil {
		t.Fatalf("NewFileItem(): %s", err)
	} else {
		if item == nil {
			t.Fatal("NewFileItem(): nil")
		}
		if got := item.String(); got != "test" {
			t.Fatalf("NewFileItem(): %s", got)
		}
	}
}

func TestFileItem(t *testing.T) {
	item, err := NewFileItem("test/test.txt")
	if err != nil {
		t.Fatalf("NewFileItem(): %s", err)
	}
	if item == nil {
		t.Fatal("NewFileItem(): nil")
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if got := item.Path(); got != filepath.Join(dir, "test/test.txt") {
		t.Fatalf("FileItem.Path(): %s", got)
	}
	if got := item.Base(); got != "test.txt" {
		t.Fatalf("FileItem.Base(): %s", got)
	}
	if got := item.Name(); got != "test" {
		t.Fatalf("FileItem.Name(): %s", got)
	}
}
