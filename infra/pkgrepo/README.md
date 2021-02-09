<!--
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

## start the dev vm
### Environment variables required (convenient to manage with envdir...)
* DOCKER_REGISTRY
* DOCKER_USERNAME
* DOCKER_PASSWORD
* MAGMA_ROOT

DOCKER_REGISTRY should point to a location with aptly images
TODO: give `docker` role the ability to build images at provision time

MAGMA_ROOT should point to the magma git repository currently being used,
defaults to current working copy


```
vagrant up aptly
```

## testing locally
```
# from magma/lte/gateway
fab dev package:vcs=git
```

```
fab test shipit
```

```
fab  promote:test,beta,0.3.73-1560277031-53f7ae53
```

