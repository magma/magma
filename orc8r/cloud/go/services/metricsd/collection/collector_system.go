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

package collection

import (
	"fmt"
	"syscall"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
)

// DiskUsageMetricCollector is a MetricCollector which return a pair of metric
// families representing the total available disk space on the system and the
// total disk space used, respectively
type DiskUsageMetricCollector struct{}

func (*DiskUsageMetricCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, fmt.Errorf("failed to collect disk usage statistics: %s", err)
	}

	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	used := all - free

	return []*io_prometheus_client.MetricFamily{
		makeTotalDiskSpaceMetric(all),
		makeUsedDiskSpaceMetric(used),
	}, nil
}

// ProcMetricsCollector is a MetricCollector which queries /proc for
// the number of open file descriptors across all processes running on the
// machine, returning a single metric family for the count.
type ProcMetricsCollector struct{}

func (s *ProcMetricsCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	procs, err := procfs.AllProcs()
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, err
	}
	totalFds := 0
	for _, proc := range procs {
		numFds, err := proc.FileDescriptorsLen()
		if err != nil {
			return []*io_prometheus_client.MetricFamily{}, err
		}
		totalFds = totalFds + numFds
	}
	return []*io_prometheus_client.MetricFamily{
		makeOpenFileDescriptorsMetric(uint64(totalFds)),
	}, nil
}

// makeTotalDiskSpaceMetric returns a prometheus MetricFamily with a single
// gauge value that indicates how much total disk space (in bytes) the current
// host has.
func makeTotalDiskSpaceMetric(availableSpaceBytes uint64) *io_prometheus_client.MetricFamily {
	name := "disk_total"
	help := "Total disk space on the machine"

	gaugeValue := float64(availableSpaceBytes)
	return MakeSingleGaugeFamily(name, help, nil, gaugeValue)
}

// makeUsedDiskSpaceMetric returns a prometheus MetricFamily with a single
// gauge value that indicates how much total disk space (in bytes) the current
// host has used.
func makeUsedDiskSpaceMetric(usedSpaceBytes uint64) *io_prometheus_client.MetricFamily {
	name := "disk_used"
	help := "Disk space used"

	gaugeValue := float64(usedSpaceBytes)
	return MakeSingleGaugeFamily(name, help, nil, gaugeValue)
}

// makeOpenFileDescriptorsMetric returns a prometheus MetricFamily with a
// single gauge value that indicates how many file descriptors are currently
// open across all processes on the current host
func makeOpenFileDescriptorsMetric(numFileDescriptors uint64) *io_prometheus_client.MetricFamily {
	name := "num_file_descriptors"
	help := "Total open file descriptors on the machine"

	gaugeValue := float64(numFileDescriptors)
	return MakeSingleGaugeFamily(name, help, nil, gaugeValue)
}
