################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

FROM ubuntu:bionic

RUN apt-get update && \
    apt-get install -y gnupg && \
    apt-get clean

RUN echo "deb http://repo.aptly.info/ squeeze main" > /etc/apt/sources.list.d/aptly.list && \
    apt-key adv --keyserver keys.gnupg.net --recv-keys ED75B5A4483DA07C && \
    apt-get update && \
    apt-get install aptly ca-certificates -y && \
    apt-get clean

ADD aptly.conf /etc/aptly.conf

# required for aptly
RUN apt-get install -y gnupg1 gpgv1


# add additional required software
RUN apt-get install -y sudo python less 2>&1 | tail

# create a normal user to manage repository
RUN groupadd -g 12345 aptly-user
RUN useradd -m -o -u 12345 -g 12345 -s /bin/bash aptly-user
RUN groupadd -g 34567 aptly
RUN gpasswd -a aptly-user aptly


RUN mkdir /aptly
RUN chown aptly-user:aptly-user /aptly

ADD insecurekeygen.txt /home/aptly-user/

RUN mkdir -p /var/run/aptly
RUN chown aptly-user:aptly /var/run/aptly

USER aptly-user
WORKDIR /home/aptly-user
