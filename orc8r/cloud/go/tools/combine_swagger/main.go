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

/*
	Package main holds the combine_swagger executable.

	combine_swagger is a custom tool to naively merge a set of Swagger YAML
	files. The functionality is similar to go-swagger's mixin command.

	The Swagger files are expected to have been assembled to an input
	directory, along with a "common" Swagger YAML file. Merges are handled
	naively, as a full overwrite, and with no guarantees on ordering -- except
	that the "common" file's subfields receive precedence over other subfields.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/obsidian/swagger/spec"
	"magma/orc8r/cloud/go/tools/combine_swagger/combine"
	"magma/orc8r/cloud/go/tools/combine_swagger/generate"

	"github.com/golang/glog"
)

func main() {
	inDir := flag.String("in", "", "Input directory")
	commonFilepath := flag.String("common", "", "Common definitions filepath")
	outFilepath := flag.String("out", "", "Output directory")
	generateStandAloneSpec := flag.Bool("standalone", true, "Generate standalone specs")

	flag.Parse()

	fmt.Printf("Reading Swagger specs from directory:\n%s\n\n", *inDir)
	fmt.Printf("Reading common spec from file:\n%s\n\n", *commonFilepath)

	yamlCommon, yamlSpecs, err := combine.Load(*commonFilepath, *inDir)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Printf("Combining specs together...\n\n")
	combined, warnings, err := spec.Combine(yamlCommon, yamlSpecs)
	if err != nil {
		glog.Fatal(err)
	}
	if warnings != nil {
		glog.Fatalf("Some Swagger spec traits were overwritten or unable to be read: %+v", warnings)
	}

	fmt.Printf("Writing combined Swagger spec to file:\n%s\n\n", *outFilepath)
	err = combine.Write(combined, *outFilepath)
	if err != nil {
		glog.Fatal(err)
	}

	if *generateStandAloneSpec {
		err := generate.GenerateStandaloneSpecs(*inDir, os.Getenv("MAGMA_ROOT"))
		if err != nil {
			glog.Fatalf("Error generating standalone Swagger specs %+v", err)
		}
	}
}
