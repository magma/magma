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

package handlers

// addadmin handler is a 'wrapper' of the add handler which creates an admin
// operator with all wildcard permissions
//
import (
	"fmt"
	"log"
	"os"
	"strings"

	"magma/orc8r/cloud/go/tools/commands"
)

func init() {
	cmd := CommandRegistry.Add(
		"add-admin",
		"Add a new Admin (Supervisor) Operator which has all permissions"+
			" for all entities and create a new certificate for the admin",
		addAdmin)

	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s [OPTIONS] <Admin ID>\n", os.Args[0], cmd.Name())
		f.PrintDefaults()
	}

	f.StringVar(&certFName, "cert", "",
		"Name used for Admin's certificate files: Name.pem, Name.key.pem. "+
			"Defaults to: '<Admin ID>_cert' if not provided")

	addInit(f) // see common_add.go
}

func addAdmin(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	oid := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(oid) == 0 {
		f.Usage()
		log.Fatalf("A single Admin Id must be provided.")
	}
	acl := CreateAdminACL()
	if len(acl) == 0 {
		panic("BROKEN CreateAdminACL()")
	}
	return addACL(oid, acl) // see common_add.go
}
