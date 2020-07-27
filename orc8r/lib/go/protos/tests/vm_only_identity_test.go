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

package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"magma/orc8r/lib/go/protos"
)

// TestHashableIdentitiesTypeNames uses Go compiler facilities to parse
// protoc generated identity.pb.go source file and verify that number
// of all types satisfying "isIdentity_Value" interface is equal to number
// of Hashable Identity types in the type name table.
// So, this test should fail if a new type added without a new entry in the
// table
func TestHashableIdentitiesTypeNames(t *testing.T) {
	genSource := "../identity.pb.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, genSource, nil, 0)
	if err != nil {
		t.Errorf("%s parse error: %s", genSource, err)
	}

	receivers := map[string]int{}
	receiversCount := 0
	for _, d := range f.Decls {
		fu, ok := d.(*ast.FuncDecl)
		if ok {
			if "isIdentity_Value" == fu.Name.Name {
				tname := "*protos." +
					fu.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Obj.Name
				receivers[tname] = receiversCount
				receiversCount++
			}
		}
	}
	typeTable := protos.GetHashableIdentitiesTable()
	if len(typeTable) != receiversCount {
		t.Errorf(
			"Number of HashableTypes %+v doesn't match number of receivers %+v",
			typeTable, receivers)
	}
	for k, v := range receivers {
		_, ok := typeTable[k]
		if !ok {
			t.Errorf(
				"Receiver [%d] %s in %+v doesn't have Hashable Identity in: %+v",
				v, k, receivers, typeTable)
		}
	}
}
