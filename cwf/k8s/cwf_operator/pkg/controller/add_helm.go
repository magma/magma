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

package controller

import (
	"time"

	helmcontroller "github.com/operator-framework/operator-sdk/pkg/helm/controller"
	"github.com/operator-framework/operator-sdk/pkg/helm/release"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	reconcilePeriod = 10 * time.Second
	charDir         = "helm-charts/cwf"
	cwfGroup        = "charts.helm.k8s.io"
	cwfVersion      = "v1alpha1"
	cwfKind         = "Cwf"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, AddHelmController)
}

// AddHelmController adds a helm controller to the manager for the cwf
// helm chart
func AddHelmController(mgr manager.Manager) error {
	cwfHelmGVK := schema.GroupVersionKind{
		Group:   cwfGroup,
		Version: cwfVersion,
		Kind:    cwfKind,
	}
	cwfHelmChartOptions := helmcontroller.WatchOptions{
		GVK:                     cwfHelmGVK,
		ManagerFactory:          release.NewManagerFactory(mgr, charDir),
		ReconcilePeriod:         reconcilePeriod,
		WatchDependentResources: true,
		OverrideValues:          map[string]string{},
	}
	return helmcontroller.Add(mgr, cwfHelmChartOptions)
}
