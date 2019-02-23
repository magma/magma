package service

import (
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	log.Printf("RPC [%s] %s", code, msg)
	return status.Errorf(code, msg)
}

func Error(code codes.Code, err error) error {
	log.Printf("RPC [%s] %s", code, err)
	return status.Error(code, err.Error())
}
