// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/facebookincubator/symphony/store/sign"
	"github.com/google/wire"
	"go.opencensus.io/plugin/ochttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Config defines the s3 signer configuration.
type Config struct {
	Bucket   string
	Region   string
	Endpoint string
	Expire   time.Duration
}

// AddFlagsVar adds the flags used by this package to the Kingpin application.
func AddFlagsVar(a *kingpin.Application, config *Config) {
	a.Flag("s3.bucket", "s3 bucket name").
		Envar("S3_BUCKET").
		Required().
		StringVar(&config.Bucket)
	a.Flag("s3.region", "s3 bucket region").
		Envar("S3_REGION").
		StringVar(&config.Region)
	a.Flag("s3.endpoint", "s3 service endpoint").
		Envar("S3_ENDPOINT").
		StringVar(&config.Endpoint)
	a.Flag("s3.expire", "s3 signature expiration").
		Envar("S3_EXPIRE").
		Default("24h").
		DurationVar(&config.Expire)
}

// AddFlags adds the flags used by this package to the Kingpin application.
func AddFlags(a *kingpin.Application) *Config {
	config := &Config{}
	AddFlagsVar(a, config)
	return config
}

// Provider is a Wire provider set that produces a signer from config.
var Provider = wire.NewSet(
	NewSigner,
	wire.Bind(new(sign.Signer), new(*Signer)),
)

// NewSigner create a new aws-s3 signer.
func NewSigner(cfg Config) (*Signer, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	if cfg.Region == "" {
		cfg.Region, err = s3manager.GetBucketRegion(
			context.Background(), sess, cfg.Bucket, endpoints.UsEast1RegionID,
		)
		if err != nil {
			return nil, fmt.Errorf("resolving bucket %q region: %w", cfg.Bucket, err)
		}
	}
	if cfg.Expire == 0 {
		cfg.Expire = 24 * time.Hour
	}

	return &Signer{
		client: s3.New(sess, &aws.Config{
			Region: aws.String(cfg.Region),
			Endpoint: func() *string {
				if cfg.Endpoint != "" {
					return aws.String(cfg.Endpoint)
				}
				return nil
			}(),
			HTTPClient: &http.Client{
				Transport: &ochttp.Transport{},
			},
		}),
		bucket: aws.String(cfg.Bucket),
		expire: cfg.Expire,
	}, nil
}

// Signer signs s3 bucket requests.
type Signer struct {
	client *s3.S3
	bucket *string
	expire time.Duration
}

func (s *Signer) Sign(ctx context.Context, op sign.Operation, key, filename string) (string, error) {
	var req *request.Request
	switch op {
	case sign.GetObject:
		req, _ = s.client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: s.bucket,
			Key:    aws.String(key),
		})
	case sign.PutObject:
		req, _ = s.client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: s.bucket,
			Key:    aws.String(key),
		})
	case sign.DeleteObject:
		req, _ = s.client.DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: s.bucket,
			Key:    aws.String(key),
		})
	case sign.DownloadObject:
		req, _ = s.client.GetObjectRequest(&s3.GetObjectInput{
			Bucket:                     s.bucket,
			Key:                        aws.String(key),
			ResponseContentDisposition: aws.String("attachment; filename=" + filename),
		})
	default:
		return "", fmt.Errorf("invalid sign operation: %d", op)
	}
	req.SetContext(ctx)

	url, _, err := req.PresignRequest(s.expire)
	if err != nil {
		return "", fmt.Errorf("pre-sign s3 request: %w", err)
	}
	return url, nil
}
