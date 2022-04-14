# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL:=/bin/bash

.PHONY: build clean clean_gen download fmt fullgen gen gen_prots lint test tidy vet order_imports

build::
	go install ./...

clean::
	go clean ./...

clean_gen::
	for f in $$(find . -name '*.pb.go' ! -path '*/migrations/*') ; do rm $$f ; done
	for f in $$(find . -name '*_swaggergen.go' ! -path '*/migrations/*') ; do rm $$f ; done

download::
	go mod download

fmt::
	gofmt -s -w .

gen::
	go generate ./...

# gen_protos generates the protos for a module.
#
# HACK: this is gross, but the best solution given the current constraints.
# Here are the issues, and how we'll progressively fix them
# 	(1) Using a for-loop. Waiting for Golang 1.16, which will allow using the
#		most recent version of protoc.
# 	(2) Overriding field_mask. This version of protoc points to the WKTs where
#		field_mask has no go_package defined. Same resolution as (1).
# 	(3) Special-case Prometheus include. For some reason we originally shimmed
#		this vendored dependency in in a hacky way, and we use it in
#		Go, Python, and C++. I spent 3 hours trying to fix its C++
#		compilation to no success, so for now this is what we get.
# 	(4) Duplicated protoc calls. Need to move all protos to a single (IDL)
#		directory per module.
gen_protos::
	cd $(MAGMA_ROOT) ; \
	for x in $$(find $(MODULE_NAME)/protos -name '*.proto') ; do \
		protoc \
			--proto_path $(MAGMA_ROOT) \
			--proto_path $(MAGMA_ROOT)/orc8r/protos/prometheus \
			--proto_path $(PROTO_INCLUDES) \
			--go_out=plugins=grpc,Mgoogle/protobuf/field_mask.proto=google.golang.org/genproto/protobuf/field_mask:$(MAGMA_ROOT)/.. \
			$${x} || exit 1 ; \
	done ; \
	for x in $$(find $(MODULE_NAME)/cloud/go -name '*.proto') ; do \
		protoc \
			--proto_path $(MAGMA_ROOT) \
			--proto_path $(MAGMA_ROOT)/orc8r/protos/prometheus \
			--proto_path $(PROTO_INCLUDES) \
			--go_opt=paths=source_relative \
			--go_out=plugins=grpc,Mgoogle/protobuf/field_mask.proto=google.golang.org/genproto/protobuf/field_mask:. \
			$${x} || exit 1 ; \
	done


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

# reorder imports with goimports
order_imports:
	for f in $$(find . -type f -name '*.go' ! -name '*.pb.go' ! -name '*_swaggergen.go') ; do $(MAGMA_ROOT)/orc8r/cloud/order_go_imports.sh $${f} || exit 1 ; done

lint:
	golangci-lint run

swagger_tools:
	go install magma/orc8r/cloud/go/tools/swaggergen

ifndef TEST_RESULTS_DIR
TEST_RESULTS_DIR := $(MAGMA_ROOT)/orc8r/cloud/test-results
export TEST_RESULTS_DIR
endif
test::
	mkdir -p $(TEST_RESULTS_DIR)
	$(eval NAME ?= $(shell pwd | tr / _))
	gotestsum --junitfile $(TEST_RESULTS_DIR)/$(NAME).xml

tidy::
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
	awk '!/\.pb\.go|_swaggergen\.go|\/mocks\/|\/tools\/|\/blobstore\//' $(COVER_FILE) > $(COVER_FILE).tmp && mv $(COVER_FILE).tmp $(COVER_FILE)
cover_pre:
	mkdir -p $(COVER_DIR)
