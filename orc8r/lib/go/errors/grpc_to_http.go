/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetHttpStatusCode(err error) int {
	if _, ok := err.(ClientInitError); ok {
		return http.StatusServiceUnavailable
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.OK:
			return http.StatusOK

		case codes.Canceled:
		case codes.Unknown:
		case codes.DeadlineExceeded:
		case codes.ResourceExhausted:
		case codes.Aborted:
		case codes.OutOfRange:
		case codes.Internal:
		case codes.DataLoss:
			return http.StatusInternalServerError

		case codes.InvalidArgument:
		case codes.AlreadyExists:
		case codes.FailedPrecondition:
			return http.StatusBadRequest

		case codes.NotFound:
			return http.StatusNotFound

		case codes.PermissionDenied:
			return http.StatusForbidden

		case codes.Unimplemented:
			return http.StatusNotImplemented

		case codes.Unavailable:
			return http.StatusServiceUnavailable

		case codes.Unauthenticated:
			return http.StatusUnauthorized

		default:
			return http.StatusInternalServerError
		}
	}

	return http.StatusInternalServerError
}
