# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG ENV=standard
FROM python:3.9.2-slim-buster AS builder

COPY dp/cloud/python/magma/db_service/requirements.txt /dp/cloud/python/magma/db_service/requirements.txt
WORKDIR /dp/cloud/python/magma/db_service/migrations
RUN pip3 install --upgrade pip --no-cache-dir -r ../requirements.txt

COPY dp/cloud/python/magma/db_service /dp/cloud/python/magma/db_service/
COPY dp/cloud/python/magma/mappings /dp/cloud/python/magma/mappings/

FROM builder AS final
ENV PYTHONPATH=/:/dp/cloud/python
ENV ALEMBIC_CONFIG=./alembic.ini

ENTRYPOINT ["python"]
CMD ["../db_initialize.py"]
