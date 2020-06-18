# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Base alpine golang image
FROM golang:1.14-alpine

# Install build tools
RUN apk add --no-cache git gcc musl-dev

# Download dependencies
WORKDIR /src/symphony
COPY go.* ./
RUN go mod download

# Download gotestsum tool
RUN go get gotest.tools/gotestsum

# Switch to intergration tests directory
WORKDIR /src/symphony/integration/tests

# Running container executes integration tests
CMD ["gotestsum", "-f", "short-verbose", "--", "-tags=integration"]
