package servicers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"magma/fbinternal/cloud/go/services/download"
	"magma/orc8r/lib/go/registry"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
	"golang.org/x/net/http2"
)

const (
	s3Prefix    = "/s3/" // This is temporary while we figure out the right conventions
	s3PrefixLen = len(s3Prefix)
)

func RootHandler(config DownloadServiceConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.V(2).Infof("Download svc called with '%s'", r.URL.Path)

		if !(strings.HasPrefix(r.URL.Path, s3Prefix)) {
			// Return 400 since the download path is unknown
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}

		appName := "default"
		imageName := r.URL.Path[s3PrefixLen:]
		paths := strings.SplitN(imageName, "/", 2)
		if len(paths) > 1 {
			appName = paths[0]
			imageName = paths[1]
		}
		appConfig, ok := config.apps[appName]
		if !ok || len(imageName) == 0 {
			// Unknown app or empty image name
			http.Error(w, "Unknown app name", http.StatusBadRequest)
			return
		}

		s3Handler(w, appConfig, imageName, r.Header.Get("Range"))
	})
}

func s3Handler(
	w http.ResponseWriter,
	appConfig AppConfig,
	imageName string,
	reqRange string,
) {
	sess, err := session.NewSession()
	if err != nil {
		glog.Errorf("Error creating AWS session: %s\n", err.Error())
		http.Error(w, "Error creating AWS session", http.StatusInternalServerError)
		return
	}

	key := appConfig.s3SubFolder + imageName
	bucket := appConfig.s3Bucket
	s3svc := s3.New(sess, &aws.Config{Region: aws.String(appConfig.s3Region)})
	s3Req := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	if reqRange != "" {
		s3Req.SetRange(reqRange)
	}
	glog.V(2).Infof("Fetching s3 object %s from bucket %s", key, bucket)
	result, err := s3svc.GetObject(&s3Req)
	if err != nil {
		glog.Errorf("Error with s3 GetObj: %s", err.Error())
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// TODO - should be status 206 for partial content
	if result.ContentRange != nil {
		w.Header().Set("Content-Range", *result.ContentRange)
	}
	// Set content length if known
	if result.ContentLength != nil {
		w.Header().Set("Content-Length", strconv.FormatInt(*result.ContentLength, 10))
	}

	_, err = io.Copy(w, result.Body)
	if err != nil {
		glog.Errorf("Copy failed: %s", err.Error())
		http.Error(w, "Copy failed", http.StatusInternalServerError)
	}
}

// Run starts the 'download' microservice
func Run() {
	config := InitServiceConfig()

	port, err := registry.GetServicePort(download.ServiceName)
	if err != nil {
		glog.Fatalf("Unable to determine port to run download service: %s", err)
	}
	glog.V(2).Infof("Listening on port %d...", port)

	// Note, the following works well, but it only supports http2, non-SSL.
	// All attempts to run a go-based server that supported both http/1.x and http2,
	// without requiring SSL, failed.
	server := http2.Server{}
	// Listen doesnt't bind to both ipv4 and ipv6 right now
	// See https://github.com/golang/go/issues/9334
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.Fatalf("net.Listen err: %v", err)
	}
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			glog.Fatalf("tcpListener.Accept err: %v", err)
		}
		server.ServeConn(conn, &http2.ServeConnOpts{
			Handler: RootHandler(config),
		})
	}
}
