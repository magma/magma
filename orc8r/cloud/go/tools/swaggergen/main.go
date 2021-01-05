/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
	Package main holds the swaggergen executable.

	swaggergen is a custom tool to generate Go code from Swagger 2.0 spec files
	in a way that allows Magma to keep Swagger files modular.
	Because one Swagger file can reference definitions from any number of other
	Swagger files across modules, we have extended the Swagger spec with some
	extra metadata

	```
	magma-gen-meta:
	  go-package: magma/orc8r/cloud/go/tools/swaggergen/testdata/importer2/models
	  dependencies:
	    - 'orc8r/cloud/go/tools/swaggergen/testdata/base.yml'
	    - 'orc8r/cloud/go/tools/swaggergen/testdata/importer.yml'
	  temp-gen-filename: importer2-swagger.yml
	  output-dir: orc8r/cloud/go/tools/swaggergen/testdata/importer2
	  types:
	    - go-struct-name: ImportingChainDef
		  filename: importing_chain_def_swaggergen.go
	```

	Think of `dependencies` as an import statement. These filepaths should be
	relative to the --root command line flag (defaults to $MAGMA_ROOT).

	`tmp-gen-filename` is what you want this file to be named when its contents
	are copied into the working directory when a dependent Swagger spec is
	being codegened from. All dependent files should reference definitions
	inside this file using this filename, as if the file was in the same
	directory, e.g.

	```
	$ref: './importer2-swagger.yml#/definitions/foo'
	```

	`output-dir` specifies where you want to generate the models to, relative
	to --root.

	During the code generation step, swaggergen will read the entire dependency
	tree for the target Swagger spec and copy all files in that tree to the
	working directory as whatever is specified in each file's
	`tmp-gen-filename`. Swaggergen will clean up these temporary files before
	exiting.

	The code generation step is configurable via a go-swagger config file.
	Ref: https://goswagger.io/use/template_layout.html

	During the code modification step, swaggergen will rewrite references in
	the generated code to structs owned by different Swagger spec files to the
	implementations in go-package specified by the Swagger spec file which
	owns the type. The struct type ownership is defined by the types key of
	the meta map.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/tools/swaggergen/generate"

	"github.com/golang/glog"
)

func main() {
	targetFilepath := flag.String("target", "", "Target swagger spec to generate code from")
	configFilepath := flag.String("config", "", "Config file for go-swagger command")
	rootDir := flag.String("root", os.Getenv("MAGMA_ROOT"), "Root path to resolve dependency and output directories based on")
	flag.Parse()

	if *targetFilepath == "" {
		glog.Fatal("Target file must be specified")
	}
	if *configFilepath == "" {
		glog.Fatal("Config file must be specified")
	}
	if *rootDir == "" {
		glog.Fatal("Root directory must be specified, or MAGMA_ROOT must be in env")
	}

	cwd, _ := os.Getwd()
	fmt.Printf("Generating swagger types for %s from directory %s\n", *targetFilepath, cwd)

	specs, err := generate.ParseSwaggerDependencyTree(*targetFilepath, *rootDir)
	if err != nil {
		glog.Fatalf("Error parsing swagger spec dependency tree: %v\n", err)
	}

	err = generate.GenerateModels(*targetFilepath, *configFilepath, *rootDir, specs)
	if err != nil {
		glog.Fatalf("Error generating default swagger models: %v\n", err)
	}

	err = generate.RewriteGeneratedRefs(*targetFilepath, *rootDir, specs)
	if err != nil {
		glog.Fatalf("Error rewriting generated swagger models: %v\n", err)
	}
}
