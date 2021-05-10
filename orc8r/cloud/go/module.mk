# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
SHELL := /bin/bash
.PHONY: build clean clean_gen download fmt gen lint test tidy vet migration_plugin

build::
	go install ./...

clean::
	go clean ./...

clean_gen:
	for f in $$(find . -name '*.pb.go' ! -path '*/migrations/*') ; do rm $$f ; done
	for f in $$(find . -name '*_swaggergen.go' ! -path '*/migrations/*') ; do rm $$f ; done

download:
	go mod download

fmt::
	gofmt -s -w .

gen::
	go generate ./...


# swagger.v1.yml files are expected to be arranged one-per-service, at the
# following location
#
#	MODULE/cloud/go/services/SERVICE/obsidian/models/swagger.v1.yml
#
# copy_swagger_files copies Swagger files to the tmp directory under the name
#
#	SERVICE.swagger.v1.yml
#
# For example
#	- Before: lte/cloud/go/services/policydb/obsidian/models/swagger.v1.yml
#	- After: orc8r/cloud/swagger/specs/partial/policydb.swagger.v1.yml
copy_swagger_files:
	for f in $$(find . -name swagger.v1.yml) ; do cp $$f $${SWAGGER_V1_PARTIAL_SPECS_DIR}/$$(echo $$f | sed -r 's/.*\/services\/([^\/]*)\/obsidian\/models\/(swagger\.v1\.yml)/\1.\2/g') ; done

lint:
	golangci-lint run

swagger_tools:
	go install magma/orc8r/cloud/go/tools/combine_swagger

ifndef TEST_RESULTS_DIR
TEST_RESULTS_DIR := $(MAGMA_ROOT)/orc8r/cloud/test-results
export TEST_RESULTS_DIR
endif
test::
	mkdir -p $(TEST_RESULTS_DIR)
	$(eval NAME ?= $(shell pwd | tr / _))
	gotestsum --junitfile $(TEST_RESULTS_DIR)/$(NAME).xml

tidy:
	go mod tidy

tools:: $(TOOL_DEPS)
$(TOOL_DEPS): %:
	go install $*

vet::
	go vet -composites=false ./...

ifndef COVER_DIR
COVER_DIR := $(MAGMA_ROOT)/orc8r/cloud/coverage
export COVER_DIR
endif
COVER_FILE=$(COVER_DIR)/$(MODULE_NAME).gocov
cover: tools cover_pre
	go-acc ./... --covermode count --output $(COVER_FILE)
	# Don't measure coverage for tools and generated files
	awk '!/\.pb\.go|_swaggergen\.go|\/mocks\/|\/tools\/|\/blobstore\/ent\//' $(COVER_FILE) > $(COVER_FILE).tmp && mv $(COVER_FILE).tmp $(COVER_FILE)
cover_pre:
	mkdir -p $(COVER_DIR)

# for configurator data migration
migration_plugin:
	if [[ -d ./tools/migrations/m003_configurator/plugin ]]; then \
		go build -buildmode=plugin -o $(PLUGIN_DIR)/migrations/m003_configurator/$(PLUGIN_NAME).so ./tools/migrations/m003_configurator/plugin; \
	fi
