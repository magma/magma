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
	"database/sql"
	"fmt"
	"sort"

	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/lte/cloud/go/services/subscriberdb/state"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type IPLookup interface {
	// Initialize the backing store.
	Initialize() error

	// GetIPs returns the set of IMSIs each IP is marked as assigned to, keyed
	// by IP address.
	GetIPs(networkID string, ips []string) ([]*protos.IPMapping, error)

	// SetIPs assigns an IP to an IMSI under a particular APN, one per mapping.
	SetIPs(networkID string, mappings []*protos.IPMapping) error
}

const (
	ipLookupTableName = "subscriberdb_ip_to_imsi"

	ipLookupNidCol  = "network_id"
	ipLookupIpCol   = "ip"
	ipLookupImsiCol = "imsi_and_apn"
)

type ipLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewIPLookup(db *sql.DB, builder sqorc.StatementBuilder) IPLookup {
	return &ipLookup{db: db, builder: builder}
}

func (l *ipLookup) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(ipLookupTableName).
			IfNotExists().
			Column(ipLookupNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(ipLookupIpCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(ipLookupImsiCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			PrimaryKey(ipLookupNidCol, ipLookupIpCol, ipLookupImsiCol).
			Unique(ipLookupNidCol, ipLookupImsiCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize IP lookup table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *ipLookup) GetIPs(networkID string, ips []string) ([]*protos.IPMapping, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(ipLookupIpCol, ipLookupImsiCol).
			From(ipLookupTableName).
			Where(squirrel.Eq{ipLookupNidCol: networkID, ipLookupIpCol: ips}).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "select IMSIs for IPs %v", ips)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetIPs")

		var mappings []*protos.IPMapping
		for rows.Next() {
			m := &protos.IPMapping{}
			imsiVal := ""
			err = rows.Scan(&m.Ip, &imsiVal)
			if err != nil {
				return nil, errors.Wrap(err, "select IMSIs for IPs, SQL row scan error")
			}
			m.Imsi, m.Apn, err = state.GetIMSIAndAPNFromMobilitydStateKey(imsiVal)
			if err != nil {
				return nil, err
			}
			mappings = append(mappings, m)
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select IMSIs for IPs, SQL rows error")
		}

		sort.Slice(mappings, func(i, j int) bool { return mappings[i].String() < mappings[j].String() })
		return mappings, nil
	}
	txRet, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.([]*protos.IPMapping)
	return ret, nil
}

func (l *ipLookup) SetIPs(networkID string, mappings []*protos.IPMapping) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetIPs")

		for _, m := range mappings {
			_, err := l.builder.
				Insert(ipLookupTableName).
				Columns(ipLookupNidCol, ipLookupIpCol, ipLookupImsiCol).
				Values(networkID, m.Ip, fmt.Sprintf("%s.%s", m.Imsi, m.Apn)).
				OnConflict(
					[]sqorc.UpsertValue{{Column: ipLookupIpCol, Value: m.Ip}},
					ipLookupNidCol, ipLookupImsiCol,
				).
				RunWith(sc).
				Exec()
			if err != nil {
				return nil, errors.Wrapf(err, "insert IP mapping %+v", m)
			}
		}

		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}
