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

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"fbc/lib/go/radius/dictionary"
	"fbc/lib/go/radius/dictionarygen"
)

type Refs map[string]string

func (r Refs) Set(v string) error {
	s := strings.Split(v, string(os.PathListSeparator))
	if len(s) != 2 {
		return errors.New("invalid format")
	}
	if _, exists := r[s[0]]; exists {
		return errors.New("type already exists")
	}
	if len(s[0]) == 0 || len(s[1]) == 0 {
		return errors.New("empty type and/or package name")
	}
	r[s[0]] = s[1]
	return nil
}

func (r Refs) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	first := true
	for typ, pkg := range r {
		if first {
			b.WriteString(", ")
			first = false
		}
		b.WriteString(typ)
		b.WriteRune(os.PathListSeparator)
		b.WriteString(pkg)
	}
	b.WriteByte('}')
	return b.String()
}

type Set map[string]struct{}

func (s Set) Set(v string) error {
	s[v] = struct{}{}
	return nil
}

func (s Set) List() []string {
	values := make([]string, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	sort.Strings(values)
	return values
}

func (s Set) String() string {
	values := s.List()
	return fmt.Sprintf("%#q", values)
}

func main() {
	refs := make(Refs)
	ignored := make(Set)
	packageName := flag.String("package", "main", "generated package name")
	outputFile := flag.String("output", "-", "output file (\"-\" writes to standard out)")
	flag.Var(&refs, "ref", `external package reference (format: "attribute`+string(os.PathListSeparator)+`package")`)
	flag.Var(ignored, "ignore", `attributes names to ignore`)
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	parser := dictionary.Parser{
		Opener: &dictionary.FileSystemOpener{},

		IgnoreIdenticalAttributes: true,
	}

	dict, err := parser.ParseFile(flag.Arg(0))
	if err != nil {
		fmt.Printf("radius-dict-gen: %s\n", err)
		os.Exit(1)
	}

	// Generate dictionary code
	g := dictionarygen.Generator{
		Package:            *packageName,
		IgnoredAttributes:  ignored.List(),
		ExternalAttributes: refs,
	}
	generated, err := g.Generate(dict)
	if err != nil {
		fmt.Printf("radius-dict-gen: %s\n", err)
		os.Exit(1)
	}

	// Write generated code to file
	if *outputFile == "-" {
		os.Stdout.Write(generated)
	} else {
		outFile, err := os.Create(*outputFile)
		if err != nil {
			fmt.Printf("radius-dict-gen: %s\n", err)
			os.Exit(1)
		}

		if _, err := outFile.Write(generated); err != nil {
			fmt.Printf("radius-dict-gen: %s\n", err)
			os.Exit(1)
		}

		if err := outFile.Close(); err != nil {
			fmt.Printf("radius-dict-gen: %s\n", err)
			os.Exit(1)
		}
	}
}
