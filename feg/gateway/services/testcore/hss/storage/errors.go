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

package storage

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnknownSubscriberError indicates that the subscriber ID was not found.
type UnknownSubscriberError struct {
	imsi string
}

func (err UnknownSubscriberError) Error() string {
	return fmt.Sprintf("Subscriber '%s' not found", err.imsi)
}

// NewUnknownSubscriberError creates an UnknownSubscriberError.
func NewUnknownSubscriberError(imsi string) UnknownSubscriberError {
	return UnknownSubscriberError{imsi: imsi}
}

// AlreadyExistsError indicates that the subscriber ID was already in the store
// and so it cannot be added again.
type AlreadyExistsError struct {
	imsi string
}

func (err AlreadyExistsError) Error() string {
	return fmt.Sprintf("Subscriber '%s' already exists", err.imsi)
}

// NewAlreadyExistsError creates an AlreadyExistsError.
func NewAlreadyExistsError(imsi string) AlreadyExistsError {
	return AlreadyExistsError{imsi: imsi}
}

// InvalidArgumentError indicates that one of the arguments given to a function
// was not valid.
type InvalidArgumentError struct {
	msg string
}

func (err InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument error: %s", err.msg)
}

// NewInvalidArgumentError creates an InvalidArgumentError.
func NewInvalidArgumentError(msg string) InvalidArgumentError {
	return InvalidArgumentError{msg: msg}
}

// ConvertStorageErrorToGrpcStatus converts any custom hss strogae error
// into a gRPC status error.
func ConvertStorageErrorToGrpcStatus(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case AlreadyExistsError:
		return status.Errorf(codes.AlreadyExists, err.Error())
	case InvalidArgumentError:
		return status.Errorf(codes.InvalidArgument, e.msg)
	case UnknownSubscriberError:
		return status.Errorf(codes.NotFound, err.Error())
	default:
		return status.Errorf(codes.Unknown, err.Error())
	}
}
