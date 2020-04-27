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
	"encoding/json"
	"errors"
	"testing"
)

// Create a client and execute test cases.
func doTestClient(t *testing.T, f func(*Client)) {
	if client := New(); client == nil {
		t.Fatal("New() return nil")
	} else {
		f(client)
	}
}

// A loader for testing.
type testLoader struct {
	items map[string]string
}

func newTestLoader(items map[string]string) Loader {
	return &testLoader{items: items}
}

func (l *testLoader) Load(target string, next Next) ([]byte, error) {
	if value, found := l.items[target]; found {
		return []byte(value), nil
	}
	return next()
}

func TestClient(t *testing.T) {
	fooLoader := newTestLoader(map[string]string{"foo": "foo"})
	barLoader := newTestLoader(map[string]string{"bar": "bar"})
	errLoader := LoaderFunc(func(string, Next) ([]byte, error) {
		return nil, errors.New("error loader")
	})

	doTestClient(t, func(client *Client) {
		if client.Use(fooLoader) == nil {
			t.Fatal("Client.Use() return nil")
		}
		if client.Use(barLoader) == nil {
			t.Fatal("Client.Use() return nil")
		}

		if data, err := client.Load("not-found"); err == nil {
			t.Fatal("No error")
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
	})

	doTestClient(t, func(client *Client) {
		if client.Use(errLoader) == nil {
			t.Fatal("Client.Use() return nil")
		}

		if data, err := client.Load("not-found"); err == nil {
			t.Fatal("No error")
		} else {
			if data != nil {
				t.Fatal(string(data))
			}
		}
	})
}

func TestClient_LoadBuffer(t *testing.T) {
	doTestClient(t, func(client *Client) {
		client.Use(newTestLoader(map[string]string{"test": "test"}))

		if buffer, err := client.LoadBuffer("test"); err != nil {
			t.Fatal(err)
		} else {
			if buffer == nil {
				t.Fatal("Nil buffer")
			}
			if s := buffer.String(); s != "test" {
				t.Fatal(s)
			}
		}
	})

	doTestClient(t, func(client *Client) {
		client.Use(LoaderFunc(func(string, Next) ([]byte, error) {
			return nil, errors.New("error loader")
		}))

		if buffer, err := client.LoadBuffer("test"); err == nil {
			t.Fatal("No error")
		} else {
			if buffer != nil {
				t.Fatal(buffer.String())
			}
		}
	})
}

// A bindable object for testing binding.
type testBindable struct {
	Key string `json:"key"`
}

func newBindable() *testBindable {
	return new(testBindable)
}

func (o *testBindable) BindFrom(data []byte) error {
	return json.Unmarshal(data, o)
}

func TestClient_LoadAndBind(t *testing.T) {
	doTestClient(t, func(client *Client) {
		client.Use(LoaderFunc(func(string, Next) ([]byte, error) {
			return nil, errors.New("error loader")
		}))

		b := newBindable()

		if err := client.LoadAndBind("test", b); err == nil {
			t.Fatal("No error")
		} else {
			if b.Key != "" {
				t.Fatal(b.Key)
			}
		}
	})

	doTestClient(t, func(client *Client) {
		client.Use(newTestLoader(map[string]string{"test": `{"key":"foo"}`}))
		b := newBindable()

		if err := client.LoadAndBind("test", b); err != nil {
			t.Fatal(err)
		} else {
			if b.Key != "foo" {
				t.Fatal(b.Key)
			}
		}
	})
}

// Test client bindings.
func doTestClientWithObject(t *testing.T, f func(*Client, interface{}) string) {
	type object struct {
		Key string `json:"key" xml:"key" toml:"key"`
	}

	doTestClient(t, func(client *Client) {
		o := new(object)
		if want := f(client, o); want != o.Key {
			t.Fatalf("Want %s, Got %s", want, o.Key)
		}
	})
}

func TestClient_LoadJSON(t *testing.T) {
	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		client.Use(newTestLoader(map[string]string{"foo": `{"key":"foo"}`}))
		if err := client.LoadJSON("foo", o); err != nil {
			t.Fatal(err)
		}
		return "foo"
	})

	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		if err := client.LoadJSON("foo", o); err == nil {
			t.Fatal("No error")
		} else {
			if err != ErrNotFound {
				t.Fatal(err)
			}
		}
		return ""
	})

	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		client.Use(newTestLoader(map[string]string{"foo": `{bar}`}))
		if err := client.LoadJSON("foo", o); err == nil {
			t.Fatal("No error")
		}
		return ""
	})
}

func TestClient_LoadXML(t *testing.T) {
	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		client.Use(newTestLoader(map[string]string{"foo": `<xml><key>foo</key></xml>`}))
		if err := client.LoadXML("foo", o); err != nil {
			t.Fatal(err)
		}
		return "foo"
	})

	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		if err := client.LoadXML("foo", o); err == nil {
			t.Fatal("No error")
		} else {
			if err != ErrNotFound {
				t.Fatal(err)
			}
		}
		return ""
	})

	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		client.Use(newTestLoader(map[string]string{"foo": `><`}))
		if err := client.LoadXML("foo", o); err == nil {
			t.Fatal("No error")
		}
		return ""
	})
}

func TestClient_LoadTOML(t *testing.T) {
	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		client.Use(newTestLoader(map[string]string{"foo": `key = "foo"`}))
		if err := client.LoadTOML("foo", o); err != nil {
			t.Fatal(err)
		}
		return "foo"
	})

	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		if err := client.LoadTOML("foo", o); err == nil {
			t.Fatal("No error")
		} else {
			if err != ErrNotFound {
				t.Fatal(err)
			}
		}
		return ""
	})

	doTestClientWithObject(t, func(client *Client, o interface{}) string {
		client.Use(newTestLoader(map[string]string{"foo": `===`}))
		if err := client.LoadTOML("foo", o); err == nil {
			t.Fatal("No error")
		}
		return ""
	})
}
