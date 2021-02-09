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

# Requirements:
* docker
* docker-compose
* web browser (nginx container currently serves at http://localhost:8880/ )
* some .deb files to serve

# Download Images and Run
From the same directory as this file:
```
docker login facebookconnectivity-orc8r-docker.jfrog.io

docker-compose pull
docker-compose up
```

# Build and Run
From the same directory as this file:
```
docker-compose up --build
```

# Access
available users
root
aptly-user
```
docker-compose exec aptly bash --login
docker-compose exec -u root aptly bash --login
```

# Keys
in order to publish an aptly repository, aptly-user needs a gpg1 key -- see aptly/insecurekeygen.txt
or ~aptly-user/insecurekeygen.txt for sample instructions


# Usage
first upload some .deb files
```
docker-compose exec aptly mkdir /home/aptly-user/upload
docker cp ~/fbsource/fbcode/magma/.cache/apt/xenial "$(docker-compose ps -q aptly)":/home/aptly-user/upload
docker-compose exec -u root aptly chown -R aptly-user:aptly-user /home/aptly-user/upload
```

create a repository
```
docker-compose exec aptly aptly repo create -architectures=amd64 example1
```

add some packages
```
docker-compose exec aptly aptly repo add example1 upload/xenial
```

create a snapshot
```
export SNAPSHOT_DATE=$(date +%Y%m%d%H%M%S)
docker-compose exec aptly aptly snapshot create snap_${SNAPSHOT_DATE} from repo example1
```

publish snapshot by name (see aptly publish switch)
```
docker-compose exec aptly aptly publish snapshot -distribution=exampledistro snap_${SNAPSHOT_DATE}
```

publish snapshot with timestamp
```
docker-compose exec aptly aptly publish snapshot -distribution=snap_${SNAPSHOT_DATE} snap_${SNAPSHOT_DATE}
```


If everything works right, you should have a valid debian repository at http://localhost:8880/



Further Reading
https://www.aptly.info/doc/
