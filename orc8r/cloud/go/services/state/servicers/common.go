package servicers

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	return status.Error(codes.Internal, e.Error())
}
