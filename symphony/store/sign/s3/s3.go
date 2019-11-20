// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"context"
	"net/http"
	"time"

	"github.com/facebookincubator/symphony/store/sign"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
)

type (
	// Config defines the s3 signer configuration.
	Config struct {
		Bucket   string        `env:"BUCKET" long:"bucket" description:"s3 bucket name" required:"true"`
		Region   string        `env:"REGION" long:"region" description:"s3 bucket region"`
		Endpoint string        `env:"ENDPOINT" long:"endpoint" description:"s3 service endpoint"`
		Expire   time.Duration `env:"EXPIRE" long:"expire" default:"24h" description:"s3 signature expiration"`
	}

	signer struct {
		*s3.S3
		bkt    *string
		expire time.Duration
	}
)

// Set is a Wire provider set that produces a signer from config.
var Set = wire.NewSet(
	NewSigner,
)

// NewSigner create a new aws-s3 signer.
func NewSigner(cfg Config) (sign.Signer, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "creating session")
	}

	if cfg.Region == "" {
		cfg.Region, err = s3manager.GetBucketRegion(context.Background(), sess, cfg.Bucket, endpoints.UsEast1RegionID)
		if err != nil {
			return nil, errors.Wrapf(err, "getting bucket region: %s", cfg.Bucket)
		}
	}
	if cfg.Expire == 0 {
		cfg.Expire = 24 * time.Hour
	}

	return &signer{
		S3: s3.New(sess, &aws.Config{
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
		bkt:    aws.String(cfg.Bucket),
		expire: cfg.Expire,
	}, nil
}

func (s *signer) Sign(ctx context.Context, op sign.Operation, key, filename string) (string, error) {
	var req *request.Request
	switch op {
	case sign.GetObject:
		req, _ = s.GetObjectRequest(&s3.GetObjectInput{
			Bucket: s.bkt,
			Key:    aws.String(key),
		})
	case sign.PutObject:
		req, _ = s.PutObjectRequest(&s3.PutObjectInput{
			Bucket: s.bkt,
			Key:    aws.String(key),
		})
	case sign.DeleteObject:
		req, _ = s.DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket: s.bkt,
			Key:    aws.String(key),
		})
	case sign.DownloadObject:
		req, _ = s.GetObjectRequest(&s3.GetObjectInput{
			Bucket:                     s.bkt,
			Key:                        aws.String(key),
			ResponseContentDisposition: aws.String("attachment; filename=" + filename),
		})
	default:
		return "", errors.Errorf("invalid sign operation: %d", op)
	}
	req.SetContext(ctx)

	url, _, err := req.PresignRequest(s.expire)
	return url, errors.Wrap(err, "pre-sign s3 request")
}
