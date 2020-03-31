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
	"path/filepath"
	"strings"
	"testing"
)

func TestFileLoader(t *testing.T) {
	base, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(base, "txt")
	if loader == nil {
		t.Fatal("NewFileLoader() return nil")
	}

	next := func() ([]byte, error) { return nil, ErrNotFound }

	if data, err := loader.Load("test", next); err != nil {
		t.Fatal(err)
	} else {
		if s := strings.TrimSpace(string(data)); s != "TEST" {
			t.Fatal(s)
		}
	}

	loader = NewFileLoader(base, "")
	if data, err := loader.Load("test", next); err == nil {
		t.Fatal("FileLoader.Load() return nil error")
	} else {
		if err != ErrNotFound {
			t.Fatal(err)
		}
		if data != nil {
			t.Fatal(string(data))
		}
	}

	loader = NewFileLoader(base, "")
	if data, err := loader.Load("", next); err == nil {
		t.Fatal("FileLoader.Load() return nil error")
	} else {
		if data != nil {
			t.Fatal(string(data))
		}
	}
}
