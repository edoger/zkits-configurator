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

func doTestMemoryLoader(t *testing.T, f func(*MemoryLoader)) {
	if loader := NewMemoryLoader(); loader == nil {
		t.Fatal("NewMemoryLoader() return nil")
	} else {
		f(loader)
	}
}

func TestMemoryLoader(t *testing.T) {
	doTestMemoryLoader(t, func(loader *MemoryLoader) {
		if got := loader.Add("test", []byte(`foo`)); !got {
			t.Fatal(got)
		}
		if got := loader.Add("test", []byte(`bar`)); got {
			t.Fatal(got)
		}

		next := func() ([]byte, error) {
			return nil, ErrNotFound
		}

		if data, err := loader.Load("test", next); err != nil {
			t.Fatal(err)
		} else {
			if s := string(data); s != "foo" {
				t.Fatal(s)
			}
		}

		loader.Set("test", []byte(`bar`))

		if data, err := loader.Load("test", next); err != nil {
			t.Fatal(err)
		} else {
			if s := string(data); s != "bar" {
				t.Fatal(s)
			}
		}

		loader.Set("test", nil)

		if data, err := loader.Load("test", next); err == nil {
			t.Fatal("No error")
		} else {
			if err != ErrNotFound {
				t.Fatal(err)
			}
			if data != nil {
				t.Fatal(string(data))
			}
		}
	})
}
