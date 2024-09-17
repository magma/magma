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

package service

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
	grpc_proto "google.golang.org/grpc/encoding/proto"
)

// logModes
type grpcLogVerbosityLevel int

const (
	GRPCLOG_DISABLED    = grpcLogVerbosityLevel(iota) // Default value
	GRPCLOG_FULL                                      // Prints all protos
	GRPCLOG_HIDEVERBOSE                               // Prints all except verboseProtos
	grpclog_end
)

var (
	// Default list of verbose protos
	verboseProtos = map[string]grpcLogVerbosityLevel{
		"protos.Void":                  GRPCLOG_HIDEVERBOSE,
		"protos.ServiceInfo":           GRPCLOG_HIDEVERBOSE,
		"protos.HealthStatus":          GRPCLOG_HIDEVERBOSE,
		"protos.MetricsContainer":      GRPCLOG_HIDEVERBOSE,
		"protos.SubmitMetricsRequest":  GRPCLOG_HIDEVERBOSE,
		"protos.SubmitMetricsResponse": GRPCLOG_HIDEVERBOSE,
		"grpc.MetricFamilies":          GRPCLOG_HIDEVERBOSE,
		"grpc.Void":                    GRPCLOG_HIDEVERBOSE,
	}
)

// logCodec is a debugging Codec implementation for protobuf.
// It'll be used if debug GRPC printout is enabled
type logCodec struct {
	protoCodec     encoding.Codec
	verbosityLevel grpcLogVerbosityLevel
}

// Marshal of GRPC Codec interface
func (lc logCodec) Marshal(v interface{}) ([]byte, error) {
	printMessage("Sending: ", v, lc.verbosityLevel)
	return lc.protoCodec.Marshal(v)
}

// Unmarshal of GRPC Codec interface
func (lc logCodec) Unmarshal(data []byte, v interface{}) error {
	err := lc.protoCodec.Unmarshal(data, v)
	printMessage("Received: ", v, lc.verbosityLevel)
	return err
}

// Name of GRPC Codec interface
func (logCodec) Name() string {
	return grpc_proto.Name
}

func printMessage(prefix string, v interface{}, verboseLevel grpcLogVerbosityLevel) {
	var payload string
	if pm, ok := v.(proto.Message); ok {
		if verboseLevel == GRPCLOG_HIDEVERBOSE && isProtoHidden(pm) {
			// do not print verbose proto
			return
		}
		var buf bytes.Buffer
		err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: "\t", OrigName: true}).Marshal(&buf, pm)
		if err == nil {
			payload = buf.String()
		} else {
			payload = fmt.Sprintf("\n\t JSON encoding error: %v; %s", err, buf.String())
		}
	} else {
		payload = fmt.Sprintf("\n\t %T is not proto.Message; %+v", v, v)
	}
	glog.Infof("%s%T: %s", prefix, v, payload)
}

func isProtoHidden(pm proto.Message) bool {
	vType := strings.Replace(fmt.Sprintf("%T", pm), "*", "", 1) // remove possible * from proto name
	if val, ok := verboseProtos[vType]; ok && val == GRPCLOG_HIDEVERBOSE {
		return true
	}
	return false
}

// registerPrintGrpcPayloadLogCodecIfRequired will get MAGMA_PRINT_GRPC_PAYLOAD from flags or
// from EnvVar and register a new logger to print GRPC message content
func registerPrintGrpcPayloadLogCodecIfRequired() {
	verbosityLevel := getVerbosityLevelFromFlagOrEnvVar()
	if verbosityLevel == GRPCLOG_DISABLED {
		return
	}
	ls := logCodec{
		protoCodec:     encoding.GetCodec(grpc_proto.Name),
		verbosityLevel: verbosityLevel,
	}
	encoding.RegisterCodec(ls)
}

// getLogVerbosityFromFlagOrEnvVar parses verbosityLevel either from flags or Env var and converts it'
// into one of grpcLogVerbosityLevel valid values. Flag value will have priority
// It returns DISABLED in case value is not right
func getVerbosityLevelFromFlagOrEnvVar() grpcLogVerbosityLevel {
	if printGrpcPayload != 0 {
		return parseVerbosityLevel(printGrpcPayload)
	}
	logModeStr := strings.TrimSpace(os.Getenv(PrintGrpcPayloadEnv))
	if logModeStr == "" {
		return GRPCLOG_DISABLED
	}
	logModeInt, err := strconv.Atoi(logModeStr)
	if err != nil {
		return GRPCLOG_DISABLED
	}
	return parseVerbosityLevel(logModeInt)
}

func parseVerbosityLevel(logMode int) grpcLogVerbosityLevel {
	if logMode < 0 || logMode >= int(grpclog_end) {
		glog.Errorf(
			"could not parse PrintGrpcPayload properly. Value %d is not validDefaulting. Setting to 0 (DISABLED)",
			logMode)
		return GRPCLOG_DISABLED
	}
	return grpcLogVerbosityLevel(logMode)
}
