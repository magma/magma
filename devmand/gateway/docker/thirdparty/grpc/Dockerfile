# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM debian:stretch AS %%PKG%%
ARG PKG_DIR=/cache/%%PKG%%
ARG PKG_REPO_DIR=/cache/%%PKG%%/repo
ARG PKG_INSTALL_DIR=/cache/%%PKG%%/install

RUN %%INSTALL%% git make %%DEPS%%

WORKDIR $PKG_DIR
RUN git clone %%URL%% $PKG_REPO_DIR

WORKDIR $PKG_REPO_DIR
RUN git checkout %%VERSION%%
RUN git submodule update --init
RUN sed -i s/-Werror//g Makefile
RUN make -j 8
RUN make install prefix=$PKG_INSTALL_DIR
