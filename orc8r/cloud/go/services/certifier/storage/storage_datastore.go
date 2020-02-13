/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package storage

import (
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/services/certifier/protos"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	// Certificate info is stored in a dedicated table used by datastore.Api.
	CertifierTableDatastore = "certificate_info_db"
)

type certifierDatastore struct {
	store datastore.Api
}

// NewCertifierDatastore returns an initialized instance of certifierDatastore as CertifierStorage.
func NewCertifierDatastore(store datastore.Api) CertifierStorage {
	return &certifierDatastore{store: store}
}

func (c *certifierDatastore) ListSerialNumbers() ([]string, error) {
	return c.store.ListKeys(CertifierTableDatastore)
}

func (c *certifierDatastore) GetCertInfo(serialNumber string) (*protos.CertificateInfo, error) {
	infos, err := c.GetManyCertInfo([]string{serialNumber})
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		return info, nil
	}
	return nil, merrors.ErrNotFound
}

func (c *certifierDatastore) GetManyCertInfo(serialNumbers []string) (map[string]*protos.CertificateInfo, error) {
	marshaledInfos, err := c.store.GetMany(CertifierTableDatastore, serialNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get many certificate info")
	}

	ret := make(map[string]*protos.CertificateInfo)
	for sn, mInfoWrapper := range marshaledInfos {
		info := &protos.CertificateInfo{}
		err = proto.Unmarshal(mInfoWrapper.Value, info)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cert info")
		}
		ret[sn] = info
	}

	return ret, nil
}

// NOTE: datastore GetAllCertInfo doesn't execute in a single commit.
func (c *certifierDatastore) GetAllCertInfo() (map[string]*protos.CertificateInfo, error) {
	serialNumbers, err := c.ListSerialNumbers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list serial numbers")
	}
	return c.GetManyCertInfo(serialNumbers)
}

func (c *certifierDatastore) PutCertInfo(serialNumber string, certInfo *protos.CertificateInfo) error {
	marshaledInfo, err := proto.Marshal(certInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cert info")
	}

	err = c.store.Put(CertifierTableDatastore, serialNumber, marshaledInfo)
	if err != nil {
		return errors.Wrap(err, "failed to put certificate info")
	}

	return nil
}

func (c *certifierDatastore) DeleteCertInfo(serialNumber string) error {
	return c.store.Delete(CertifierTableDatastore, serialNumber)
}
