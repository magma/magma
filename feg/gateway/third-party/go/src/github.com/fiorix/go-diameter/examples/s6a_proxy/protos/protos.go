//go:generate protoc --go_out=plugins=grpc:. ./s6a.proto

// protos package encapsulates protoc generated go files for S6a GRPC Proxy.
// Use `go generate github.com/fiorix/go-diameter/v4/examples/s6a_proxy/protos` to re-generate
// protos
// As noted in:
//   https://docs.google.com/document/d/1V03LUfjSADDooDMhe-_K59EgpTEm3V8uvQRuNMAEnjg/edit#heading=h.tksmbbpjl4ya
//   "...First, go generate is intended§ to be run by the author of a package,
//    not the client of it. The author of the package generates the required Go files
//    and includes them in the package; the client does a regular go get or go build.
//    Generation through go generate is not part of the build, just a tool for package
//    authors. This avoids complicating the dependency analysis done by Go build..."
package protos
