/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// package testutil provides utilities for integration tests
package testlib

import (
	"io/ioutil"

	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
)

var sshKey = flag.String("ssh-key", "", "SSH key to connect to the gateways")
var sshUser = flag.String("ssh-user", "ubuntu", "SSH user")
var sshPort = flag.Uint("ssh-port", 22, "SSH port")

func keyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func runRemoteCommand(hostname string, bastionIP string, cmds []string) ([]string, error) {
	sshConfig := &ssh.ClientConfig{
		User:            *sshUser,
		Auth:            []ssh.AuthMethod{keyFile(*sshKey)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	bastionConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", bastionIP, *sshPort), sshConfig)
	if err != nil {
		fmt.Printf("Error %v", err)
		return nil, err
	}
	defer bastionConn.Close()

	tunnelConn, err := bastionConn.Dial("tcp", fmt.Sprintf("%s:%d", hostname, *sshPort))
	if err != nil {
		fmt.Printf("Error %v", err)
		return nil, err
	}
	defer tunnelConn.Close()

	newClientConn, channels, sshRequest, err := ssh.NewClientConn(tunnelConn, hostname, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to establish new ssh client conn [%s]: %v", hostname, err)
	}
	defer newClientConn.Close()

	conn := ssh.NewClient(newClientConn, channels, sshRequest)
	defer conn.Close()

	output := []string{}
	for _, cmd := range cmds {
		session, err := conn.NewSession()
		if err != nil {
			return nil, err
		}
		defer session.Close()

		out, err := session.CombinedOutput(cmd)
		if err != nil {
			return nil, err
		}
		output = append(output, string(out))
	}
	return output, nil
}
