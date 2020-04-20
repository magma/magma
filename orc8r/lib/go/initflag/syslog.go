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

package initflag

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"
)

const syslogFlag = "stderr2syslog"

type syslogTarget struct {
	network, addr string
	isSet         bool
}

var (
	syslogDest     syslogTarget
	stdoutToStderr = flag.Bool("stdout2stderr", false, "Redirect stdout to stderr")
	stderr, stdout *os.File
)

// syslogTarget's flag Var String implementation
func (v *syslogTarget) String() string {
	if v != nil {
		return v.network + "::" + v.addr
	}
	return ""
}

// syslogTarget's flag Var Set implementation
func (v *syslogTarget) Set(s string) error {
	if v == nil {
		return os.ErrInvalid
	}
	v.isSet = true
	v.addr, v.network = "", ""
	if len(s) > 0 && s != "true" { // "true" is a stub value for the flag without any value
		parts := strings.Split(s, "::")
		switch len(parts) {
		case 1:
			v.addr = strings.TrimSpace(parts[0])
			if len(v.addr) > 0 {
				v.network = "unixgram"
			}
		case 2:
			v.network = strings.TrimSpace(parts[0])
			v.addr = strings.TrimSpace(parts[1])
		default:
			v.isSet = false
			return os.ErrInvalid
		}
	}
	return nil
}

// IsSet returns true if the flag was set
func (v *syslogTarget) IsSet() bool {
	return (v != nil) && v.isSet
}

// Allow empty values
func (v *syslogTarget) IsBoolFlag() bool {
	return true
}

func redirectToSyslog() error {
	syslogWriter, err := syslog.Dial(syslogDest.network, syslogDest.addr, syslog.LOG_DAEMON, filepath.Base(os.Args[0]))
	if err != nil {
		return fmt.Errorf("syslog.Dial failed for %s, error: %v", syslogDest.String(), err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		return err
	}
	stderr, os.Stderr = os.Stderr, writer
	log.SetOutput(writer) // if anyone uses std log pkg

	go func() {
		scanner := bufio.NewScanner(reader)
		err = nil
		for err == nil && scanner.Scan() {
			err = writeToSyslog(syslogWriter, scanner.Text())
			if err != nil {
				break
			}
		}
		var msg string
		// Unexpected error, reset stderr, log error and return
		if err != nil {
			msg = fmt.Sprintf("syslog write error: %v\n", err)
		} else {
			msg = fmt.Sprintf("unexpected end of stderr stream error: %v\n", scanner.Err())
		}
		syslogWriter.Emerg(msg)
		os.Stderr = stderr
		log.SetOutput(stderr)
		if *stdoutToStderr && stdout != nil {
			os.Stdout = stdout
		}
		log.Print(msg) // also log into stderr
		writer.Close()
		reader.Close()
		syslogWriter.Close()
	}()
	return nil
}

func writeToSyslog(syslogWriter *syslog.Writer, msg string) error {
	switch messageSeverity(msg) {
	case 'F': // Fatal
		return syslogWriter.Emerg(msg)
	case 'A': // Alert
		return syslogWriter.Alert(msg)
	case 'C': // Critical
		return syslogWriter.Crit(msg)
	case 'E': // Error
		return syslogWriter.Err(msg)
	case 'W': // Warning
		return syslogWriter.Warning(msg)
	case 'I': // Info
		return syslogWriter.Info(msg)
	case 'D': // Debug
		return syslogWriter.Debug(msg)
	}
	return syslogWriter.Notice(msg) // Notice otherwise
}

func messageSeverity(msg string) uint8 {
	// check if it's a glog message starting with W0102 ...
	if len(msg) > 6 && msg[0] >= 'D' && msg[0] <= 'W' && msg[5] == ' ' && areDigits(msg[1:5]) {
		return msg[0]
	}
	return 'N' // Notice by default
}

func areDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
