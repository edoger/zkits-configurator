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
	"strings"
	"testing"
)

func doTestFileLoader(t *testing.T, f func(*FileLoader)) {
	if loader := NewFileLoader(); loader == nil {
		t.Fatal("NewFileLoader() return nil")
	} else {
		f(loader)
	}
}

func TestFileLoader(t *testing.T) {
	doTestFileLoader(t, func(loader *FileLoader) {
		if err := loader.AddDir("test", ".txt"); err != nil {
			t.Fatal(err)
		}

		next := func() ([]byte, error) {
			return nil, ErrNotFound
		}

		if data, err := loader.Load("test", next); err != nil {
			t.Fatal(err)
		} else {
			if s := strings.TrimSpace(string(data)); s != "TEST" {
				t.Fatal(s)
			}
		}

		if data, err := loader.Load("unknown", next); err == nil {
			t.Fatal("No error")
		} else {
			if err != ErrNotFound {
				t.Fatal(err)
			}
			if data != nil {
				t.Fatal(string(data))
			}
		}

		loader.files["test"] = "test/sub"

		if data, err := loader.Load("test", next); err == nil {
			t.Fatal("No error")
		} else {
			if data != nil {
				t.Fatal(string(data))
			}
		}
	})

	doTestFileLoader(t, func(loader *FileLoader) {
		if err := loader.AddDir("test/not-found"); err == nil {
			t.Fatal("No error")
		}
		if err := loader.AddDir("\n"); err == nil {
			t.Fatal("No error")
		}
	})

	doTestFileLoader(t, func(loader *FileLoader) {
		if err := loader.AddDir("test", ".txt"); err != nil {
			t.Fatal(err)
		}
		if err := loader.AddDir("test", ".json"); err == nil {
			t.Fatal("No error")
		}
	})

	doTestFileLoader(t, func(loader *FileLoader) {
		if err := loader.AddDir("test", ".txt", ".json"); err == nil {
			t.Fatal("No error")
		}
	})

	doTestFileLoader(t, func(loader *FileLoader) {
		if err := loader.AddDir("test", ".txt"); err != nil {
			t.Fatal(err)
		}
		if err := loader.AddDir("test", ".txt"); err != nil {
			t.Fatal(err)
		}
	})

	doTestFileLoader(t, func(loader *FileLoader) {
		if err := loader.AddFile("test/sub/foo.txt"); err != nil {
			t.Fatal(err)
		}
		if err := loader.AddFile("test/test.txt"); err != nil {
			t.Fatal(err)
		}
		if err := loader.AddFile("test/test.json"); err == nil {
			t.Fatal("No error")
		}
		if err := loader.AddFile("test/sub"); err == nil {
			t.Fatal("No error")
		}
	})
}
