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

func TestFileLoader(t *testing.T) {
	loader := NewFileLoader("test")
	if loader == nil {
		t.Fatal("NewFileLoader() return nil")
	}
	if err := loader.Initialize(); err != nil {
		t.Fatal(err)
	}
	if err := loader.Initialize(); err != nil {
		t.Fatal(err)
	}

	next := func() ([]byte, error) { return nil, ErrNotFound }

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

	if err := NewFileLoader("test/error").Initialize(); err == nil {
		t.Fatal("No error")
	}
}
