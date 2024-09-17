//go:build tools
// +build tools

/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

// Put all binary tool dependencies in here so they can be tracked by the go
// module.
// dependencies listed here should also be listed in TOOL_DEPS in ./Makefile
import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/ory/go-acc"
	_ "github.com/vektra/mockery/v2"
	_ "github.com/wadey/gocovmerge"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
)
