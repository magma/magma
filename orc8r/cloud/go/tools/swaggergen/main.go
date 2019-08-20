/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

/*

swaggergen is a custom tool to generate Go code from swagger 2.0 spec files
in a way that allows Magma to keep swagger files modular.
Because one swagger file can reference definitions from any number of other
swagger files across modules, we have extended the swagger spec with some
extra metadata:

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

Think of `dependencies` as an `import` statement. These filepaths should be
relative to the --root command line flag (defaults to $MAGMA_ROOT).

`temp-gen-filename` is what you want this file to be named when its contents
are copied into the working directory when a dependent swagger spec is being
codegened from. All dependent files should reference definitions inside this
file using this filename, as if the file was in the same directory, e.g.:

```
  $ref: './importer2-swagger.yml#/definitions/foo'
```

`output-dir` specifies where you want to generate the models to, relative
to --root.

During the code generation step, swaggergen will read the entire dependency
tree for the target swagger spec and copy all files in that tree to the
working directory as whatever is specified in each file's `temp-gen-filename`.
swaggergen will clean up these temporary files before exiting.

During the code modification step, swaggergen will rewrite references in the
generated code to structs owned by different swagger spec files to the
implementations in `go-package` specified by the swagger spec file which
owns the type. The struct type ownership is defined by the `types` key of
the meta map.

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"magma/orc8r/cloud/go/tools/swaggergen/generate"
)

func main() {
	targetFile := flag.String("target", "", "Target swagger spec to generate code from")
	templateFile := flag.String("template", "", "swagger template file for code generation")
	rootDir := flag.String("root", os.Getenv("MAGMA_ROOT"), "Root path to resolve dependency and output directories based on")
	flag.Parse()

	if *targetFile == "" {
		log.Fatal("target file must be specified")
	}
	if *templateFile == "" {
		log.Fatal("template file must be specified")
	}
	if *rootDir == "" {
		log.Fatal("root dir must be specified, or MAGMA_ROOT has to be in env")
	}

	fmt.Printf("Generating swagger types for %s\n", *targetFile)
	err := generate.GenerateModels(*targetFile, *templateFile, *rootDir)
	if err != nil {
		log.Fatalf("Failed to generate swagger models: %v\n", err)
	}

	fmt.Printf("Rewriting generated swagger types for %s\n", *targetFile)
	err = generate.RewriteGeneratedRefs(*targetFile, *rootDir)
	if err != nil {
		log.Fatalf("Failed to rewrite generated swagger models: %v\n", err)
	}
}
