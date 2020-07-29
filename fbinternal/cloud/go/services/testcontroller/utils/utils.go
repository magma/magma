/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
)

func IsNewerVersion(v1 string, v2 string) (bool, error) {
	ver1, err := convertVersionStringToStruct(v1)
	if err != nil {
		return false, err
	}
	ver2, err := convertVersionStringToStruct(v2)
	if err != nil {
		return false, err
	}
	return ver2.timestamp.After(ver1.timestamp) && lessThanOrEqualTo(ver1.version, ver2.version), nil
}

func lessThanOrEqualTo(v1 *semver.Version, v2 *semver.Version) bool {
	return v1.LessThan(*v2) || v1.Equal(*v2)
}

func convertVersionStringToStruct(version string) (packageVersion, error) {
	v := strings.Split(version, "-")
	if len(v) != 3 {
		return packageVersion{}, errors.Errorf("package version string '%s' not in form 'version-timestamp-hash'", version)
	}

	i, err := strconv.ParseInt(v[1], 10, 64)
	if err != nil {
		return packageVersion{}, errors.Wrap(err, "unable to parse package timestamp")
	}

	return packageVersion{
		version:   semver.New(v[0]),
		timestamp: time.Unix(i, 0),
		hash:      v[2],
	}, nil
}

type packageVersion struct {
	version   *semver.Version
	timestamp time.Time
	hash      string
}
