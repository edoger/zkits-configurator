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

// The configuration loader.
type Loader interface {
	Load(string, func() ([]byte, error)) ([]byte, error)
}

// The configuration loader function.
type LoaderFunc func(string, func() ([]byte, error)) ([]byte, error)

// The loader can call the given trigger to give the load permission to
// the next available loader.
func (f LoaderFunc) Load(target string, next func() ([]byte, error)) ([]byte, error) {
	return f(target, next)
}
